package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	secretKey = "your-secret-key"
)

var (
	ipFamily           = flag.Int("ip-family", 6, "6 to IPv6, 4 to IPv4")
	merchantIP         = flag.String("merchant-ip", "::1", "Merchant Service IP Address")
	merchantPort       = flag.Int("merchant-port", 50051, "Merchant Service Port")
	rateLimiterIP      = flag.String("rate-limiter-ip", "::1", "Rate Limiter Service IP Address")
	rateLimiterPort    = flag.Int("rate-limiter-port", 50052, "Rate Limiter Service Port")
	ledgerIP           = flag.String("ledger-ip", "::1", "Ledger Service IP Address")
	ledgerPort         = flag.Int("ledger-port", 50053, "Ledger Service Port")
	bankIP             = flag.String("bank-ip", "127.0.0.1", "Bank IP Address")
	bankPort           = flag.Int("bank-port", 8000, "Bank Port")
	merchantAddress    = ""
	rateLimiterAddress = ""
	ledgerAddress      = ""
	bankAddress        = ""
)

func main() {

	flag.Parse()
	if ipFamily == nil || merchantIP == nil || merchantPort == nil {
		log.Fatal("requires ip-family, merchant-ip, and merchant-port")
	}

	switch *ipFamily {
	case 6:
		merchantAddress = fmt.Sprintf("[%s]:%d", *merchantIP, *merchantPort)
		rateLimiterAddress = fmt.Sprintf("[%s]:%d", *rateLimiterIP, *rateLimiterPort)
		ledgerAddress = fmt.Sprintf("[%s]:%d", *ledgerIP, *ledgerPort)
	case 4:
		merchantAddress = fmt.Sprintf("%s:%d", *merchantIP, *merchantPort)
		rateLimiterAddress = fmt.Sprintf("%s:%d", *rateLimiterIP, *rateLimiterPort)
		ledgerAddress = fmt.Sprintf("%s:%d", *ledgerIP, *ledgerPort)

	default:
		log.Fatalf("invalid ip version: %d", *ipFamily)
	}

	// always IPv4
	bankAddress = fmt.Sprintf("%s:%d", *bankIP, *bankPort)
	if !strings.HasPrefix(bankAddress, "http://") {
		bankAddress = fmt.Sprintf("http://%s", bankAddress)
	}

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type", "Accept", "Access-Control-Allow-Origin"}

	router.Use(cors.New(config))

	router.GET("/login", loginHandler)

	router.POST("/payment", authMiddleware, rateLimitMiddleware, createPaymentHandler)
	router.PUT("/payment", restrictMiddleware, updatePaymentHandler)
	router.GET("/payment/:id", authMiddleware, rateLimitMiddleware, readPaymentHandler)

	router.Run(":8080")
}
