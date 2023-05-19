package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentStatus int

const (
	Created PaymentStatus = iota
	Pending
	Success
	Fail
)

func (ps PaymentStatus) String() string {
	switch ps {
	case 0:
		return "CREATED"
	case 1:
		return "PENDING"
	case 2:
		return "SUCCESS"
	case 3:
		return "FAIL"
	default:
		return fmt.Sprintf("%d", ps)
	}
}

type Payment struct {
	ID               uuid.UUID  `json:"id"`
	MerchantID       uuid.UUID  `json:"merchant_id"`
	Amount           float64    `json:"amount"`
	Currency         string     `json:"currency"`
	PurchaseTime     time.Time  `json:"purchase_time"`
	ValidationMethod string     `json:"validation_method"`
	Card             CreditCard `json:"card"`
	Metadata         string     `json:"metadata"`
	Status           string     `json:"status"`
	BankPaymentID    uuid.UUID  `json:"bank_payment_id"`
	BankRequestTime  time.Time  `json:"bank_request_time"`
	BankResponseTime time.Time  `json:"bank_response_time"`
	BankMessage      string     `json:"bank_message"`
}

func (p Payment) GetPurchaseTimeStr() string {
	return p.PurchaseTime.Format("2006-01-02T15:04:05.000")
}

func (p Payment) GetBankRequestTimeStr() string {
	return p.BankRequestTime.Format("2006-01-02T15:04:05.000")
}

func (p Payment) GetBankResponseTimeStr() string {
	return p.BankResponseTime.Format("2006-01-02T15:04:05.000")
}
