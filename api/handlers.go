package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thiagolcmelo/payment-gateway/api/bank"
	"github.com/thiagolcmelo/payment-gateway/api/entities"
	"github.com/thiagolcmelo/payment-gateway/api/ledger"
	"github.com/thiagolcmelo/payment-gateway/api/merchant"
)

func loginHandler(c *gin.Context) {
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		log.Printf("username=%s, password=%s", username, password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing credentials"})
		return
	}

	ms := merchant.NewMerchantService(c, merchantAddress)

	merchantID, ok := ms.Validate(username, password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create a new JWT token
	claims := MerchantClaims{
		Username: username,
		ID:       merchantID,
		Exp:      time.Now().Add(time.Hour * 24),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Return the token to the client
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func createPaymentHandler(c *gin.Context) {
	var body createPaymentRequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("could not parse request: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := body.validate(); err != nil {
		log.Printf("could not validate request: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := c.MustGet("claims").(MerchantClaims)
	ms := merchant.NewMerchantService(c, merchantAddress)
	m, err := ms.Get(claims.ID)
	if err != nil {
		log.Printf("could not retrieve merchant: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	p := entities.Payment{
		MerchantID:       m.ID,
		Amount:           body.Amount,
		Currency:         body.Currency,
		PurchaseTime:     body.getPurchaseTime(),
		ValidationMethod: body.ValidationMethod,
		Card:             body.Card,
		Metadata:         body.Metadata,
		Status:           fmt.Sprint(entities.Created),
	}

	ls := ledger.NewLedgerService(c, ledgerAddress)
	p, err = ls.CreatePayment(p)
	if err != nil {
		log.Printf("could not create payment: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bs := bank.NewBankService(c, bankAddress)
	p, err = bs.RelayPaymentRequest(m, p)
	if err != nil {
		log.Printf("could not relay payment to bank: %v", err)
		p, err = ls.SetPaymentFail(p)
		if err != nil {
			log.Printf("could not set payment status to fail in the ledger: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		p, err = ls.SetPaymentPending(p)
		if err != nil {
			log.Printf("could not set payment status to pending in the ledger: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"id": p.ID.String(), "status": p.Status, "bank_message": p.BankMessage})
}

func updatePaymentHandler(c *gin.Context) {
	var body bankMessage
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("could not parse bank message: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "acknowledge": false})
		return
	}

	if err := body.validate(); err != nil {
		log.Printf("invalid bank message: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "acknowledge": false})
		return
	}

	bankPaymentID, err := uuid.Parse(body.ID)
	if err != nil {
		log.Printf("could not parse bank payment id: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid bank payment id", "acknowledge": false})
		return
	}

	ls := ledger.NewLedgerService(c, ledgerAddress)
	p, err := ls.ReadPaymentUsingBankReference(bankPaymentID)
	if err != nil {
		log.Printf("could not find payment: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid bank payment id", "acknowledge": false})
		return
	}
	p.BankResponseTime = time.Now()
	p.BankMessage = body.Message

	if body.Success {
		_, err = ls.SetPaymentSuccess(p)
		if err != nil {
			log.Printf("could not set payment to success: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error", "acknowledge": false})
			return
		}
	} else {
		_, err = ls.SetPaymentFail(p)
		if err != nil {
			log.Printf("could not set payment to fail: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error", "acknowledge": false})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"acknowledge": true})
}

func readPaymentHandler(c *gin.Context) {
	// Parse and check payment ID
	pID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("could not parse payment id: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	// Read payment from Ledger
	ls := ledger.NewLedgerService(c, ledgerAddress)
	p, err := ls.ReadPayment(pID)
	if err != nil {
		log.Printf("could not read payment: %v", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid payment id"})
		return
	}

	// Get Info about merchant requesting payment
	claims := c.MustGet("claims").(MerchantClaims)
	merchantId := claims.ID

	// Assert mechant is owner of payment
	if p.MerchantID != merchantId {
		log.Printf("merchant %s trying to read unauthorized payment id: %s", merchantId.String(), pID.String())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid payment id"})
		return
	}

	c.JSON(http.StatusOK, p)
}

type createPaymentRequestBody struct {
	Amount           float64             `json:"amount"`
	Currency         string              `json:"currency"`
	PurchaseTime     string              `json:"purchate_time"`
	ValidationMethod string              `json:"validation_method"`
	Card             entities.CreditCard `json:"card"`
	Metadata         string              `json:"metadata"`
}

func (p *createPaymentRequestBody) getPurchaseTime() time.Time {
	purchaseTime, err := time.Parse("2006-01-02T15:04:05.000", p.PurchaseTime)
	if err != nil {
		return time.Now()
	}
	return purchaseTime
}

func (p *createPaymentRequestBody) validate() error {
	_, err := time.Parse("2006-01-02T15:04:05.000", p.PurchaseTime)
	if err != nil {
		return err
	}
	if p.Amount < 0 {
		return fmt.Errorf("invalid amount: %f", p.Amount)
	}
	if p.Currency == "" {
		return fmt.Errorf("invalid currency")
	}
	if p.ValidationMethod == "" {
		return fmt.Errorf("invalid validation method")
	}
	if p.Metadata == "" {
		return fmt.Errorf("invalid metadata")
	}
	return nil
}

type bankMessage struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (ms bankMessage) validate() error {
	if ms.ID == "" {
		return fmt.Errorf("missing id")
	}
	if ms.Message == "" {
		return fmt.Errorf("missing message")
	}
	return nil
}
