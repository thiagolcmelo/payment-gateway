package ratelimiter

import (
	"context"
	"log"

	"github.com/google/uuid"
	rpcRateLimiter "github.com/thiagolcmelo/payment-gateway/ratelimiter/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RateLimiterService struct {
	ctx      context.Context
	address  string
	failOpen bool
}

func NewRateLimiterService(ctx context.Context, address string, failOpen bool) *RateLimiterService {
	return &RateLimiterService{
		ctx:      ctx,
		address:  address,
		failOpen: failOpen,
	}
}

func (rls *RateLimiterService) Allow(id uuid.UUID) bool {
	// fail open
	failOpen := true

	conn, err := grpc.Dial(rls.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("rate limiter service is unreachable at address: %s, %v", rls.address, err)
		return failOpen
	}
	defer conn.Close()

	rateLimiterClient := rpcRateLimiter.NewRateLimiterServiceClient(conn)

	req := &rpcRateLimiter.AllowRequest{
		Id: id.String(),
	}
	resp, err := rateLimiterClient.Allow(rls.ctx, req)
	if err != nil {
		log.Printf("error rate limiting merchant: %v", err)
		return failOpen
	}
	return resp.Allow
}
