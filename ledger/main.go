package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
	"github.com/thiagolcmelo/payment-gateway/ledger/pb"
	"github.com/thiagolcmelo/payment-gateway/ledger/storage"
	"github.com/thiagolcmelo/payment-gateway/ledger/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	portFlag      = flag.Int("port", 50053, "The server port")
	hostFlag      = flag.String("host", "0.0.0.0", "The server host")
	ipVersionFlag = flag.Int("ip-version", 4, "The server ip version (4 for IPv4, 6 for IPv6)")
)

type server struct {
	storage storage.Storage
	pb.UnimplementedLedgerServiceServer
}

func newServerWithMemoryStorage() *server {
	return &server{
		storage: memory.NewMemoryStorage(),
	}
}

func (s *server) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	card, err := entity.NewCreditCard(req.Card.Number, req.Card.Name, int(req.Card.ExpireMonth), int(req.Card.ExpireYear), int(req.Card.Cvv))
	if err != nil {
		log.Printf("error parsing credit card in CreatePayment: %v", err)
		return nil, err
	}

	payment, err := entity.NewPayment(req.MerchantId, float64(req.Amount), req.Currency, req.PurchaseTimeUtc, req.ValidationMethod, card, req.Metadata)
	if err != nil {
		log.Printf("error parsing payment in CreatePayment: %v", err)
		return nil, err
	}
	payment.Status = entity.Created

	id, err := s.storage.Create(payment)
	if err != nil {
		log.Printf("error saving payment in CreatePayment: %v", err)
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		Id: id.String(),
	}, nil
}

func (s *server) ReadPayment(ctx context.Context, req *pb.ReadPaymentRequest) (*pb.ReadPaymentResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in ReadPayment: %v", err)
		return nil, err
	}

	payment, err := s.storage.Read(id)
	if err != nil {
		log.Printf("error reading payment in ReadPayment: %v", err)
		return nil, err
	}

	return &pb.ReadPaymentResponse{
		Payment: &pb.Payment{
			Id:               payment.ID.String(),
			MerchantId:       payment.MerchantID.String(),
			Amount:           float32(payment.Amount), // TODO: find better solution
			Currency:         payment.Currency,
			PurchaseTimeUtc:  payment.GetPurchaseTimeStr(),
			ValidationMethod: payment.ValidationMethod,
			Card: &pb.CreditCard{
				Number:      payment.Card.Number,
				Name:        payment.Card.Name,
				ExpireMonth: int32(payment.Card.ExpireMonth),
				ExpireYear:  int32(payment.Card.ExpireYear),
				Cvv:         int32(payment.Card.CVV),
			},
			Metadata:            payment.Metadata,
			Status:              pb.PaymentStatus(payment.Status),
			BankPaymentId:       payment.BankPaymentID.String(),
			BankRequestTimeUtc:  payment.GetBankRequestTimeStr(),
			BankResponseTimeUtc: payment.GetBankResponseTimeStr(),
			BankMessage:         payment.BankMessage,
		},
	}, nil
}

func (s *server) ReadPaymentUsingBankReference(ctx context.Context, req *pb.ReadPaymentUsingBankReferenceRequest) (*pb.ReadPaymentUsingBankReferenceResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in ReadPayment: %v", err)
		return nil, err
	}

	payment, err := s.storage.ReadUsingBankReference(id)
	if err != nil {
		log.Printf("error reading payment in ReadPayment: %v", err)
		return nil, err
	}

	return &pb.ReadPaymentUsingBankReferenceResponse{
		Payment: &pb.Payment{
			Id:               payment.ID.String(),
			MerchantId:       payment.MerchantID.String(),
			Amount:           float32(payment.Amount), // TODO: find better solution
			Currency:         payment.Currency,
			PurchaseTimeUtc:  payment.GetPurchaseTimeStr(),
			ValidationMethod: payment.ValidationMethod,
			Card: &pb.CreditCard{
				Number:      payment.Card.Number,
				Name:        payment.Card.Name,
				ExpireMonth: int32(payment.Card.ExpireMonth),
				ExpireYear:  int32(payment.Card.ExpireYear),
				Cvv:         int32(payment.Card.CVV),
			},
			Metadata:            payment.Metadata,
			Status:              pb.PaymentStatus(payment.Status),
			BankPaymentId:       payment.BankPaymentID.String(),
			BankRequestTimeUtc:  payment.GetBankRequestTimeStr(),
			BankResponseTimeUtc: payment.GetBankResponseTimeStr(),
			BankMessage:         payment.BankMessage,
		},
	}, nil
}

func (s *server) UpdatePaymentToPending(ctx context.Context, req *pb.UpdatePaymentToPendingRequest) (*pb.UpdatePaymentToPendingResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in UpdatePaymentToPending: %v", err)
		return nil, err
	}

	bankPaymentID, err := uuid.Parse(req.BankPaymentId)
	if err != nil {
		log.Printf("error parsing bank uuid in UpdatePaymentToPending: %v", err)
		return nil, err
	}

	payment, err := s.storage.Read(id)
	if err != nil {
		log.Printf("error reading payment in UpdatePaymentToPending: %v", err)
		return nil, err
	}

	err = payment.SetBankRequestTimeFromStr(req.BankRequestTimeUtc)
	if err != nil {
		log.Printf("error parsing date in UpdatePaymentToPending: %v", err)
		return nil, err
	}
	payment.BankPaymentID = bankPaymentID
	payment.Status = entity.Pending

	err = s.storage.Update(payment)
	if err != nil {
		log.Printf("error updating payment in UpdatePaymentToPending: %v", err)
		return nil, err
	}

	return &pb.UpdatePaymentToPendingResponse{}, nil
}

func (s *server) UpdatePaymentToSuccess(ctx context.Context, req *pb.UpdatePaymentToSuccessRequest) (*pb.UpdatePaymentToSuccessResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	bankPaymentID, err := uuid.Parse(req.BankPaymentId)
	if err != nil {
		log.Printf("error parsing bank uuid in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	payment, err := s.storage.Read(id)
	if err != nil {
		log.Printf("error reading payment in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}
	err = payment.SetBankResponseTimeFromStr(req.BankResponseTimeUtc)
	if err != nil {
		log.Printf("error parsing date in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}
	payment.BankPaymentID = bankPaymentID
	payment.BankMessage = req.BankMessage
	payment.Status = entity.Success

	err = s.storage.Update(payment)
	if err != nil {
		log.Printf("error updating payment in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	return &pb.UpdatePaymentToSuccessResponse{}, nil
}

func (s *server) UpdatePaymentToFail(ctx context.Context, req *pb.UpdatePaymentToFailRequest) (*pb.UpdatePaymentToFailResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	var bankPaymentID uuid.UUID
	if req.BankPaymentId != nil {
		bankPaymentID, err = uuid.Parse(*req.BankPaymentId)
		if err != nil {
			log.Printf("error parsing bank uuid in UpdatePaymentToSuccess: %v", err)
			return nil, err
		}
	}

	var bankMessage string
	if req.BankMessage != nil {
		bankMessage = *req.BankMessage
	}

	payment, err := s.storage.Read(id)
	if err != nil {
		log.Printf("error reading payment in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	var bankResponseTimeUtc string
	if req.BankResponseTimeUtc != nil {
		bankResponseTimeUtc = *req.BankResponseTimeUtc
	}
	err = payment.SetBankResponseTimeFromStr(bankResponseTimeUtc)
	if err != nil {
		log.Printf("error parsing date in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}
	payment.BankPaymentID = bankPaymentID
	payment.BankMessage = bankMessage
	payment.Status = entity.Fail

	err = s.storage.Update(payment)
	if err != nil {
		log.Printf("error updating payment in UpdatePaymentToSuccess: %v", err)
		return nil, err
	}

	return &pb.UpdatePaymentToFailResponse{}, nil
}

func main() {
	var host string
	var port int
	var ipVersion int

	// prefer environment variables over flags
	envHost := os.Getenv("SERVICE_HOST")
	envPort := os.Getenv("SERVICE_PORT")
	envIpVersion := os.Getenv("SERVICE_IP_VERSION")
	if envHost != "" && envPort != "" && envIpVersion != "" {
		host = envHost
		port, _ = strconv.Atoi(envPort)
		ipVersion, _ = strconv.Atoi(envIpVersion)
	} else {
		flag.Parse()
		host = *hostFlag
		port = *portFlag
		ipVersion = *ipVersionFlag
	}

	if ipVersion == 6 {
		host = fmt.Sprintf("[%s]", host)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLedgerServiceServer(s, newServerWithMemoryStorage())
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
