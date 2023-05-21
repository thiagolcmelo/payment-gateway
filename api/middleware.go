package main

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thiagolcmelo/payment-gateway/api/ratelimiter"
)

type MerchantClaims struct {
	Username string    `json:"username"`
	ID       uuid.UUID `json:"id"`
	Exp      time.Time `json:"exp"`
	jwt.StandardClaims
}

func (mc MerchantClaims) valid() bool {
	if mc.ID == uuid.Nil {
		return false
	}

	if mc.Exp.Before(time.Now()) {
		return false
	}

	return true
}

func authMiddleware(c *gin.Context) {
	var tokenString string
	authorizationHeader := c.GetHeader("Authorization")
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) == 2 {
		tokenString = strings.TrimSpace(parts[1])
	} else {
		tokenString = strings.TrimSpace(authorizationHeader)
	}

	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
		return
	}

	var claims MerchantClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("err: %v", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not parse token"})
		return
	}

	if claims.valid() && token.Valid {
		c.Set("claims", claims)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
	}
}

func rateLimitMiddleware(c *gin.Context) {
	claims := c.MustGet("claims").(MerchantClaims)
	rls := ratelimiter.NewRateLimiterService(c, rateLimiterAddress, true)
	if !rls.Allow(claims.ID) {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
	} else {
		c.Next()
	}
}

func restrictMiddleware(ipToEnable string) func(c *gin.Context) {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			log.Printf("error getting client ip: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Printf("ip=%s, expected=%s", ip, ipToEnable)
		if ip != ipToEnable {
			log.Printf("unauthorized ip: %s", ip)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}

	// allowedNetwork := "10.123.123.0/30" // Specify the network address and mask
	// // Get the client's IP address
	// ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	// if err != nil {
	// 	log.Printf("error getting client ip: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }
	// // Check if the client's IP is within the allowed network
	// _, allowedCIDR, err := net.ParseCIDR(allowedNetwork)
	// if err != nil {
	// 	log.Printf("error verifying ip: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }
	// clientIP := net.ParseIP(ip)
	// if allowedCIDR.Contains(clientIP) {
	// 	// IP is allowed, proceed to the next handler
	// 	c.Next()
	// } else {
	// 	// IP is not allowed, return a Forbidden response
	// 	c.AbortWithStatus(http.StatusForbidden)
	// }
}
