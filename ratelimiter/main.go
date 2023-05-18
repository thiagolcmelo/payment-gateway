package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	merchant "github.com/thiagolcmelo/payment-gateway/merchant/pb"
	"github.com/thiagolcmelo/payment-gateway/ratelimiter/pb"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port         = flag.Int("port", 50052, "The server port")
	ipFamily     = flag.Int("ip-family", 6, "6 to IPv6, 4 to IPv4")
	merchantIP   = flag.String("merchant-ip", "::1", "Merchant Service IP Address")
	merchantPort = flag.Int("merchant-port", 50051, "Merchant Service Port")
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type server struct {
	merchantServiceAddress string
	clients                map[uuid.UUID]*client
	pb.UnimplementedRateLimiterServiceServer
	sync.Mutex
}

func newServerWithMemoryLimiter(merchantServiceAddress string) *server {
	return &server{
		merchantServiceAddress: merchantServiceAddress,
		clients:                make(map[uuid.UUID]*client),
	}
}

func (s *server) getMaxQPS(ctx context.Context, id uuid.UUID) (int, error) {
	conn, err := grpc.Dial(s.merchantServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 0, fmt.Errorf("merchant service is unreachable at address: %s", s.merchantServiceAddress)
	}
	defer conn.Close()
	merchantClient := merchant.NewMerchantServiceClient(conn)

	merchantReq := &merchant.GetQPSRequest{
		Id: id.String(),
	}
	merchantResp, err := merchantClient.GetQPS(context.Background(), merchantReq)
	if err != nil {
		return 0, fmt.Errorf("error reading max qps: %v", err)
	}
	return int(merchantResp.MaxQps), nil
}

func (s *server) allowClient(ctx context.Context, id uuid.UUID) (bool, error) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.clients[id]; !ok {
		maxQPS, err := s.getMaxQPS(ctx, id)
		if err != nil {
			log.Printf("could not read from merchant service: %v", err)
			return false, err
		}
		if maxQPS > 0 {
			s.clients[id] = &client{limiter: rate.NewLimiter(rate.Limit(maxQPS), 10)}
		} else {
			s.clients[id] = &client{limiter: &rate.Limiter{}}
		}

	}

	s.clients[id].lastSeen = time.Now()
	return s.clients[id].limiter.Allow(), nil
}

func (s *server) Allow(ctx context.Context, req *pb.AllowRequest) (*pb.AllowResponse, error) {
	// do not propagate errors
	id, err := uuid.Parse(req.Id)
	if err != nil {
		log.Printf("error parsing uuid in Allow: %v", err)
		return nil, err
	}

	allow, err := s.allowClient(ctx, id)
	if err != nil {
		log.Printf("could not rate limit: %v", err)
	}
	return &pb.AllowResponse{
		Allow: allow,
	}, err
}

func main() {
	flag.Parse()

	if port == nil || ipFamily == nil || merchantIP == nil || merchantPort == nil {
		log.Fatal("requires port, ip-family, merchant-ip, and merchant-port")
	}

	// Set up a connection to the gRPC server
	var merchantAddress string
	switch *ipFamily {
	case 6:
		merchantAddress = fmt.Sprintf("[%s]:%d", *merchantIP, *merchantPort)
	case 4:
		merchantAddress = fmt.Sprintf("%s:%d", *merchantIP, *merchantPort)
	default:
		log.Fatalf("invalid ip version: %d", *ipFamily)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRateLimiterServiceServer(s, newServerWithMemoryLimiter(merchantAddress))
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
