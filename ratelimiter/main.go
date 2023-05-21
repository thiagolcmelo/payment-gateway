package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
	portFlag         = flag.Int("port", 50052, "The server port")
	hostFlag         = flag.String("host", "0.0.0.0", "The server host")
	ipVersionFlag    = flag.Int("ip-version", 4, "The server ip version (4 for IPv4, 6 for IPv6)")
	merchantHostFlag = flag.String("merchant-host", "0.0.0.0", "Merchant Service host")
	merchantPortFlag = flag.Int("merchant-port", 50051, "Merchant Service port")
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
	merchantResp, err := merchantClient.GetQPS(ctx, merchantReq)
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

func getEnvOrFlag[T int | string](env string, flagVal *T, conv func(string) (T, error)) T {
	if envVal := os.Getenv(env); envVal != "" {
		v, err := conv(envVal)
		if err != nil {
			return *flagVal
		}
		return v
	}
	return *flagVal
}

func main() {
	var (
		network      string = "tcp4"
		ipVersion    int    = getEnvOrFlag("IP_VERSION", ipVersionFlag, strconv.Atoi)
		host         string = getEnvOrFlag("RATE_LIMITER_SERVICE_HOST", hostFlag, func(v string) (string, error) { return v, nil })
		port         int    = getEnvOrFlag("RATE_LIMITER_SERVICE_PORT", portFlag, strconv.Atoi)
		merchantHost string = getEnvOrFlag("MERCHANT_SERVICE_HOST", merchantHostFlag, func(v string) (string, error) { return v, nil })
		merchantPort int    = getEnvOrFlag("MERCHANT_SERVICE_PORT", merchantPortFlag, strconv.Atoi)
	)

	if ipVersion == 6 {
		host = fmt.Sprintf("[%s]", host)
		merchantHost = fmt.Sprintf("[%s]", merchantHost)
		network = "tcp6"
	}

	address := fmt.Sprintf("%s:%d", host, port)
	merchantAddress := fmt.Sprintf("%s:%d", merchantHost, merchantPort)

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRateLimiterServiceServer(s, newServerWithMemoryLimiter(merchantAddress))
	reflection.Register(s)
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
