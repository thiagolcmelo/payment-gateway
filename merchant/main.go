package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/merchant/entities"
	"github.com/thiagolcmelo/payment-gateway/merchant/pb"
	"github.com/thiagolcmelo/payment-gateway/merchant/storage"
	"github.com/thiagolcmelo/payment-gateway/merchant/storage/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func generateMemoryStorage() *memory.Storage {
	ms := memory.NewMemoryStorage()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s/data/merchants.json", dir))
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var merchants []entity.Merchant
	err = json.Unmarshal(data, &merchants)
	if err != nil {
		log.Fatal(err)
	}

	for _, merchant := range merchants {
		_, err := ms.CreateMerchant(merchant)
		if err != nil {
			log.Fatal(err)
		}
	}

	return ms
}

type server struct {
	storage storage.Storage
	pb.UnimplementedMerchantServiceServer
}

func newServerWithMemoryStorage() *server {
	ms := generateMemoryStorage()

	return &server{
		storage: ms,
	}
}

func (s *server) GetMerchant(ctx context.Context, req *pb.GetMerchantRequest) (*pb.GetMerchantResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in GetMerchant: %v", err)
		return nil, err
	}

	merchant, err := s.storage.ReadMerchant(id)
	if err != nil {
		log.Printf("error reading storage in GetMerchant: %v", err)
		return nil, err
	}

	return &pb.GetMerchantResponse{
		Id:       merchant.ID.String(),
		Username: merchant.Username,
		Password: merchant.Password,
		Active:   merchant.Active,
		MaxQps:   int32(merchant.MaxQPS),
	}, nil
}

func (s *server) GetQPS(ctx context.Context, req *pb.GetQPSRequest) (*pb.GetQPSResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in GetQPS: %v", err)
		return nil, err
	}

	merchant, err := s.storage.ReadMerchant(id)
	if err != nil {
		log.Printf("error reading storage in GetQPS: %v", err)
		return nil, err
	}

	return &pb.GetQPSResponse{
		MaxQps: int32(merchant.MaxQPS),
	}, nil
}

func (s *server) MerchantActive(ctx context.Context, req *pb.MerchantActiveRequest) (*pb.MerchantActiveResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in MerchantActive: %v", err)
		return nil, err
	}

	merchant, err := s.storage.ReadMerchant(id)
	if err != nil {
		log.Printf("error reading storage in MerchantActive: %v", err)
		return nil, err
	}

	return &pb.MerchantActiveResponse{
		Active: merchant.Active,
	}, nil
}

func (s *server) FindMerchant(ctx context.Context, req *pb.FindMerchantRequest) (*pb.FindMerchantResponse, error) {
	id, err := s.storage.FindMerchantID(req.Username, req.Password)
	if err != nil {
		log.Printf("error checking is merchant exists: %v", err)
		return &pb.FindMerchantResponse{
			Exists: false,
			Id:     nil,
		}, nil
	}
	idStr := id.String()
	return &pb.FindMerchantResponse{
		Exists: true,
		Id:     &idStr,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMerchantServiceServer(s, newServerWithMemoryStorage())
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
