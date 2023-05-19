package merchant

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/thiagolcmelo/payment-gateway/api/entities"
	rpcMerchant "github.com/thiagolcmelo/payment-gateway/merchant/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MerchantService struct {
	ctx     context.Context
	address string
}

func NewMerchantService(ctx context.Context, address string) *MerchantService {
	return &MerchantService{
		ctx:     ctx,
		address: address,
	}
}

func (ms *MerchantService) Validate(username, password string) (uuid.UUID, bool) {
	conn, err := grpc.Dial(ms.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("merchant service is unreachable at address: %s, %v", ms.address, err)
		return uuid.Nil, false
	}
	defer conn.Close()

	merchantClient := rpcMerchant.NewMerchantServiceClient(conn)

	req := &rpcMerchant.FindMerchantRequest{
		Username: username,
		Password: password,
	}

	resp, err := merchantClient.FindMerchant(ms.ctx, req)
	if err != nil {
		log.Printf("error finding merchant: %v", err)
		return uuid.Nil, false
	}

	if !resp.Exists {
		return uuid.Nil, false
	} else if resp.Id == nil {
		log.Print("missing merchant id for existing merchant")
		return uuid.Nil, false
	}

	id, err := uuid.Parse(*resp.Id)
	if err != nil {
		log.Printf("error parsing merchant id: %v", err)
		return uuid.Nil, false
	}

	return id, true
}

func (ms *MerchantService) Get(id uuid.UUID) (entities.Merchant, error) {
	conn, err := grpc.Dial(ms.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("merchant service is unreachable at address: %s, %v", ms.address, err)
		return entities.Merchant{}, err
	}
	defer conn.Close()

	merchantClient := rpcMerchant.NewMerchantServiceClient(conn)

	req := &rpcMerchant.GetMerchantRequest{
		Id: id.String(),
	}
	resp, err := merchantClient.GetMerchant(ms.ctx, req)
	if err != nil {
		log.Printf("error getting merchant: %v", err)
		return entities.Merchant{}, err
	}

	return entities.Merchant{
		ID:       id,
		Username: resp.Username,
		Password: resp.Password,
		Name:     resp.Name,
		Active:   resp.Active,
		MaxQPS:   int(resp.MaxQps),
	}, nil
}
