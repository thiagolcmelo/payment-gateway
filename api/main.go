package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	secretKey = "your-secret-key"
)

type Merchant struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	router := gin.Default()

	router.GET("/login", loginHandler)

	router.POST("/payment", authMiddleware, createPaymentHandler)
	router.PUT("/payment", restrictMiddleware, updatePaymentHandler)
	router.GET("/payment", authMiddleware, readPaymentHandler)

	router.Run(":8080")
}

func loginHandler(c *gin.Context) {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing credentials"})
		return
	}

	// For simplicity, let's assume a hardcoded username and password.
	if username == "admin" && password == "password" {
		// Create a new JWT token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set token claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = username
		claims["id"] = username
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time

		// Generate encoded token
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		// Return the token to the client
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}
}

func authMiddleware(c *gin.Context) {

}

func restrictMiddleware(c *gin.Context) {
	allowedNetwork := "10.123.123.0/30" // Specify the network address and mask

	// Get the client's IP address
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		log.Println("error getting client ip:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Check if the client's IP is within the allowed network
	_, allowedCIDR, err := net.ParseCIDR(allowedNetwork)
	if err != nil {
		log.Println("error verifying ip:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	clientIP := net.ParseIP(ip)
	if allowedCIDR.Contains(clientIP) {
		// IP is allowed, proceed to the next handler
		c.Next()
	} else {
		// IP is not allowed, return a Forbidden response
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func createPaymentHandler(c *gin.Context) {

}

func updatePaymentHandler(c *gin.Context) {

}

func readPaymentHandler(c *gin.Context) {

}
