package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/thiagolcmelo/payment-gateway/api/entities"
)

type BankService struct {
	ctx     context.Context
	address string
}

func NewBankService(ctx context.Context, address string) *BankService {
	return &BankService{
		ctx:     ctx,
		address: address,
	}
}

func (bs *BankService) RelayPaymentRequest(m entities.Merchant, p entities.Payment) (entities.Payment, error) {
	type messageRequest struct {
		Amount           float64             `json:"amount"`
		Currency         string              `json:"currency"`
		PurchaseTime     string              `json:"purchase_time"`
		ValidationMethod string              `json:"validation_method"`
		Card             entities.CreditCard `json:"card"`
		Merchant         string              `json:"merchant"`
	}

	type messageResponse struct {
		Id      string `json:"id"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	// Create bank request payload
	payload := messageRequest{
		Amount:           p.Amount,
		Currency:         p.Currency,
		PurchaseTime:     p.GetPurchaseTimeStr(),
		ValidationMethod: p.ValidationMethod,
		Card:             p.Card,
		Merchant:         m.Name,
	}

	// Marshal payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("could not marshal json: %v", err)
		return p, err
	}

	// Create a POST request with the JSON payload
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/payment", bs.address), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("error creating request: %v", err)
		return p, err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending request: %v", err)
		return p, err
	}
	defer resp.Body.Close()
	p.BankRequestTime = time.Now()

	// Check the response
	if resp.StatusCode != http.StatusCreated {
		log.Printf("error relaying payment to bank, status code: %v (%s)", resp.StatusCode, resp.Status)
		// for bad request, it should read the message
		if resp.StatusCode != http.StatusBadRequest {
			return p, err
		}
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not read response: %v", err)
		return p, err
	}

	// Unmarchal response
	var responseData messageResponse
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Printf("could not unmarshal response: %v", err)
		return p, err
	}
	p.BankMessage = responseData.Message
	if !responseData.Success {
		return p, fmt.Errorf("request to bank resulted in: %s", responseData.Message)
	}

	// Assert reference id is uuid
	id, err := uuid.Parse(responseData.Id)
	if err != nil {
		log.Printf("could not parse uuid: %v", err)
		return p, err
	}
	p.BankPaymentID = id

	return p, nil
}
