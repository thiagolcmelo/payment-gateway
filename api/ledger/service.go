package ledger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thiagolcmelo/payment-gateway/api/entities"
	rpcLedger "github.com/thiagolcmelo/payment-gateway/ledger/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LedgerService struct {
	ctx     context.Context
	address string
}

func NewLedgerService(ctx context.Context, address string) *LedgerService {
	return &LedgerService{
		ctx:     ctx,
		address: address,
	}
}

func (ls *LedgerService) CreatePayment(p entities.Payment) (entities.Payment, error) {
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return p, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	req := &rpcLedger.CreatePaymentRequest{
		MerchantId:       p.MerchantID.String(),
		Amount:           float32(p.Amount),
		Currency:         p.Currency,
		PurchaseTimeUtc:  p.GetPurchaseTimeStr(),
		ValidationMethod: p.ValidationMethod,
		Card: &rpcLedger.CreditCard{
			Number:      p.Card.Number,
			Name:        p.Card.Name,
			ExpireMonth: int32(p.Card.ExpireMonth),
			ExpireYear:  int32(p.Card.ExpireYear),
			Cvv:         int32(p.Card.CVV),
		},
		Metadata: p.Metadata,
	}

	resp, err := ledgerClient.CreatePayment(ls.ctx, req)
	if err != nil {
		log.Printf("error creating payment: %v", err)
		return p, err
	}

	id, err := uuid.Parse(resp.Id)
	if err != nil {
		log.Printf("error parsing uuid: %v", err)
		return p, err
	}
	p.ID = id

	return p, nil
}

func (ls *LedgerService) SetPaymentPending(p entities.Payment) (entities.Payment, error) {
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return p, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	req := &rpcLedger.UpdatePaymentToPendingRequest{
		Id:                 p.ID.String(),
		BankPaymentId:      p.BankPaymentID.String(),
		BankRequestTimeUtc: p.GetBankRequestTimeStr(),
	}

	_, err = ledgerClient.UpdatePaymentToPending(ls.ctx, req)
	if err != nil {
		log.Printf("error updating payment to pending: %v", err)
		return p, err
	}
	p.Status = fmt.Sprint(entities.Pending)
	return p, nil
}

func (ls *LedgerService) SetPaymentSuccess(p entities.Payment) (entities.Payment, error) {
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return p, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	req := &rpcLedger.UpdatePaymentToSuccessRequest{
		Id:                  p.ID.String(),
		BankPaymentId:       p.BankPaymentID.String(),
		BankResponseTimeUtc: p.GetBankResponseTimeStr(),
		BankMessage:         p.BankMessage,
	}

	_, err = ledgerClient.UpdatePaymentToSuccess(ls.ctx, req)
	if err != nil {
		log.Printf("error updating payment to success: %v", err)
		return p, err
	}
	p.Status = fmt.Sprint(entities.Success)
	return p, nil
}

func (ls *LedgerService) SetPaymentFail(p entities.Payment) (entities.Payment, error) {
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return p, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	idStr := p.BankPaymentID.String()
	respTimeStr := p.GetBankResponseTimeStr()
	req := &rpcLedger.UpdatePaymentToFailRequest{
		Id:                  p.ID.String(),
		BankPaymentId:       &idStr,
		BankResponseTimeUtc: &respTimeStr,
		BankMessage:         &p.BankMessage,
	}

	_, err = ledgerClient.UpdatePaymentToFail(ls.ctx, req)
	if err != nil {
		log.Printf("error updating payment to fail: %v", err)
		return p, err
	}
	p.Status = fmt.Sprint(entities.Fail)
	return p, nil
}

func (ls *LedgerService) ReadPayment(id uuid.UUID) (entities.Payment, error) {
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return entities.Payment{}, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	req := &rpcLedger.ReadPaymentRequest{
		Id: id.String(),
	}

	resp, err := ledgerClient.ReadPayment(ls.ctx, req)
	if err != nil {
		log.Printf("error reading payment: %v", err)
		return entities.Payment{}, err
	}

	merchantID, err := uuid.Parse(resp.Payment.MerchantId)
	if err != nil {
		log.Printf("error parsing merchant uuid: %v", err)
	}
	bankPaymentID, err := uuid.Parse(resp.Payment.BankPaymentId)
	if err != nil {
		log.Printf("error parsing bank payment uuid: %v", err)
	}
	purchaseTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.PurchaseTimeUtc)
	if err != nil {
		log.Printf("error parsing purchate time: %v", err)
	}
	bankRequestTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.BankRequestTimeUtc)
	if err != nil {
		log.Printf("error parsing bank request time: %v", err)
	}
	bankResponseTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.BankResponseTimeUtc)
	if err != nil {
		log.Printf("error parsing bank response time: %v", err)
	}
	card := entities.CreditCard{
		Number:      resp.Payment.Card.Number,
		Name:        resp.Payment.Card.Name,
		ExpireMonth: int(resp.Payment.Card.ExpireMonth),
		ExpireYear:  int(resp.Payment.Card.ExpireYear),
		CVV:         int(resp.Payment.Card.Cvv),
	}

	return entities.Payment{
		ID:               id,
		MerchantID:       merchantID,
		Amount:           float64(resp.Payment.Amount),
		Currency:         resp.Payment.Currency,
		PurchaseTime:     purchaseTimeUTC,
		ValidationMethod: resp.Payment.ValidationMethod,
		Card:             card,
		Metadata:         resp.Payment.Metadata,
		Status:           fmt.Sprint(entities.PaymentStatus(resp.Payment.Status)),
		BankPaymentID:    bankPaymentID,
		BankRequestTime:  bankRequestTimeUTC,
		BankResponseTime: bankResponseTimeUTC,
		BankMessage:      resp.Payment.BankMessage,
	}, nil
}

func (ls *LedgerService) ReadPaymentUsingBankReference(bankPaymentID uuid.UUID) (entities.Payment, error) {
	// return entities.Payment{}, nil
	conn, err := grpc.Dial(ls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("ledger service is unreachable at address: %s, %v", ls.address, err)
		return entities.Payment{}, err
	}
	defer conn.Close()

	ledgerClient := rpcLedger.NewLedgerServiceClient(conn)

	req := &rpcLedger.ReadPaymentUsingBankReferenceRequest{
		Id: bankPaymentID.String(),
	}

	resp, err := ledgerClient.ReadPaymentUsingBankReference(ls.ctx, req)
	if err != nil {
		log.Printf("error reading payment: %v", err)
		return entities.Payment{}, err
	}

	merchantID, err := uuid.Parse(resp.Payment.MerchantId)
	if err != nil {
		log.Printf("error parsing merchant uuid: %v", err)
	}
	paymentID, err := uuid.Parse(resp.Payment.Id)
	if err != nil {
		log.Printf("error parsing payment uuid: %v", err)
	}
	purchaseTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.PurchaseTimeUtc)
	if err != nil {
		log.Printf("error parsing purchate time: %v", err)
	}
	bankRequestTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.BankRequestTimeUtc)
	if err != nil {
		log.Printf("error parsing bank request time: %v", err)
	}
	bankResponseTimeUTC, err := time.Parse("2006-01-02T15:04:05.000", resp.Payment.BankResponseTimeUtc)
	if err != nil {
		log.Printf("error parsing bank response time: %v", err)
	}
	card := entities.CreditCard{
		Number:      resp.Payment.Card.Number,
		Name:        resp.Payment.Card.Name,
		ExpireMonth: int(resp.Payment.Card.ExpireMonth),
		ExpireYear:  int(resp.Payment.Card.ExpireYear),
		CVV:         int(resp.Payment.Card.Cvv),
	}

	return entities.Payment{
		ID:               paymentID,
		MerchantID:       merchantID,
		Amount:           float64(resp.Payment.Amount),
		Currency:         resp.Payment.Currency,
		PurchaseTime:     purchaseTimeUTC,
		ValidationMethod: resp.Payment.ValidationMethod,
		Card:             card,
		Metadata:         resp.Payment.Metadata,
		Status:           fmt.Sprint(entities.PaymentStatus(resp.Payment.Status)),
		BankPaymentID:    bankPaymentID,
		BankRequestTime:  bankRequestTimeUTC,
		BankResponseTime: bankResponseTimeUTC,
		BankMessage:      resp.Payment.BankMessage,
	}, nil
}
