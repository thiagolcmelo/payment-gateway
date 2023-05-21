package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	secretKey = "your-secret-key"
)

var (
	ipVersionFlag       = flag.Int("ip-family", 6, "6 to IPv6, 4 to IPv4")
	hostFlag            = flag.String("host", "0.0.0.0", "Service host address")
	portFlag            = flag.Int("port", 8080, "Service Port")
	merchantHostFlag    = flag.String("merchant-host", "0.0.0.0", "Merchant Service host address")
	merchantPortFlag    = flag.Int("merchant-port", 50051, "Merchant Service Port")
	rateLimiterHostFlag = flag.String("rate-limiter-host", "0.0.0.0", "Rate Limiter Service host address")
	rateLimiterPortFlag = flag.Int("rate-limiter-port", 50052, "Rate Limiter Service Port")
	ledgerHostFlag      = flag.String("ledger-host", "0.0.0.0", "Ledger Service host address")
	ledgerPortFlag      = flag.Int("ledger-port", 50053, "Ledger Service Port")
	bankHostFlag        = flag.String("bank-host", "0.0.0.0", "Bank host address")
	bankPortFlag        = flag.Int("bank-port", 8000, "Bank Port")
	merchantAddress     string
	rateLimiterAddress  string
	ledgerAddress       string
	bankAddress         string
)

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
	flag.Parse()
	dummyFunc := func(v string) (string, error) { return v, nil }

	var (
		address         string
		ipVersion       int    = getEnvOrFlag("IP_VERSION", ipVersionFlag, strconv.Atoi)
		host            string = getEnvOrFlag("PAYMENT_API_SERVICE_HOST", hostFlag, dummyFunc)
		port            int    = getEnvOrFlag("PAYMENT_API_SERVICE_PORT", portFlag, strconv.Atoi)
		merchantHost    string = getEnvOrFlag("MERCHANT_SERVICE_HOST", merchantHostFlag, dummyFunc)
		merchantPort    int    = getEnvOrFlag("MERCHANT_SERVICE_PORT", merchantPortFlag, strconv.Atoi)
		rateLimiterHost string = getEnvOrFlag("RATE_LIMITER_SERVICE_HOST", rateLimiterHostFlag, dummyFunc)
		rateLimiterPort int    = getEnvOrFlag("RATE_LIMITER_SERVICE_PORT", rateLimiterPortFlag, strconv.Atoi)
		ledgerHost      string = getEnvOrFlag("LEDGER_SERVICE_HOST", ledgerHostFlag, dummyFunc)
		ledgerPort      int    = getEnvOrFlag("LEDGER_SERVICE_PORT", ledgerPortFlag, strconv.Atoi)
		bankHost        string = getEnvOrFlag("BANK_SIMULATOR_HOST", bankHostFlag, dummyFunc)
		bankPort        int    = getEnvOrFlag("BANK_SIMULATOR_PORT", bankPortFlag, strconv.Atoi)
	)

	if ipVersion == 6 {
		host = fmt.Sprintf("[%s]", host)
		merchantHost = fmt.Sprintf("[%s]", merchantHost)
		rateLimiterHost = fmt.Sprintf("[%s]", rateLimiterHost)
		ledgerHost = fmt.Sprintf("[%s]", ledgerHost)
		bankHost = fmt.Sprintf("[%s]", bankHost)
	}

	address = fmt.Sprintf("%s:%d", host, port)
	merchantAddress = fmt.Sprintf("%s:%d", merchantHost, merchantPort)
	rateLimiterAddress = fmt.Sprintf("%s:%d", rateLimiterHost, rateLimiterPort)
	ledgerAddress = fmt.Sprintf("%s:%d", ledgerHost, ledgerPort)
	bankAddress = fmt.Sprintf("http://%s:%d", bankHost, bankPort)

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

	router.Run(address)
}
