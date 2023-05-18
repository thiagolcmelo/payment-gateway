package entity

import (
	"errors"
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

var (
	// ErrNegativeAmount must be use when validating payment and amount is negative
	ErrNegativeAmount = errors.New("negative amount")
	// ErrMissingCurrency must be use when validating payment and currency is missing
	ErrMissingCurrency = errors.New("missing currency")
	// ErrMissingValidationMethod must be use when validating payment and validation method is missing
	ErrMissingValidationMethod = errors.New("missing validation method")
)

type Payment struct {
	ID               uuid.UUID
	MerchantID       uuid.UUID
	Amount           float64
	Currency         string
	PurchaseTime     time.Time
	ValidationMethod string
	Card             CreditCard
	Metadata         string
	Status           PaymentStatus
	BankPaymentID    uuid.UUID
	BankRequestTime  time.Time
	BankResponseTime time.Time
	BankMessage      string
}

// NewPayment is a factory for a payment just created
func NewPayment(
	merchantID string,
	amount float64,
	currency string,
	purchaseTimeUTC string,
	validationMethod string,
	card CreditCard,
	metadata string,
) (Payment, error) {
	purchaseTime, err := time.Parse("2006-01-02T15:04:05.000", purchaseTimeUTC)
	if err != nil {
		return Payment{}, err
	}

	merchantUUID, err := uuid.Parse(merchantID)
	if err != nil {
		return Payment{}, err
	}

	p := Payment{
		MerchantID:       merchantUUID,
		Amount:           amount,
		Currency:         currency,
		PurchaseTime:     purchaseTime,
		ValidationMethod: validationMethod,
		Card:             card,
		Metadata:         metadata,
		Status:           Created,
	}

	return p, p.Validate()
}

// Validate runs some checks to asser a payment is valid
func (p Payment) Validate() error {
	if p.Amount < 0.0 {
		return ErrNegativeAmount
	}
	if p.Currency == "" {
		return ErrMissingCurrency
	}
	if p.ValidationMethod == "" {
		return ErrMissingValidationMethod
	}
	return p.Card.Validate()
}

// Equal compares a payment with another ignoring ID only
func (p Payment) Equal(other Payment) bool {
	return p.MerchantID == other.MerchantID &&
		p.Amount == other.Amount &&
		p.Currency == other.Currency &&
		p.PurchaseTime == other.PurchaseTime &&
		p.ValidationMethod == other.ValidationMethod &&
		p.Card == other.Card &&
		p.Metadata == other.Metadata &&
		p.Status == other.Status &&
		p.BankPaymentID == other.BankPaymentID &&
		p.BankRequestTime == other.BankRequestTime &&
		p.BankResponseTime == other.BankResponseTime &&
		p.BankMessage == other.BankMessage
}

func (p *Payment) SetPurchaseTimeFromStr(value string) error {
	purchaseTime, err := time.Parse("2006-01-02T15:04:05.000", value)
	if err != nil {
		return err
	}
	p.PurchaseTime = purchaseTime
	return nil
}

func (p *Payment) GetPurchaseTimeStr() string {
	return p.PurchaseTime.Format("2006-01-02T15:04:05.000")
}

func (p *Payment) SetBankRequestTimeFromStr(value string) error {
	bankRequestTime, err := time.Parse("2006-01-02T15:04:05.000", value)
	if err != nil {
		return err
	}
	p.BankRequestTime = bankRequestTime
	return nil
}

func (p *Payment) GetBankRequestTimeStr() string {
	return p.BankRequestTime.Format("2006-01-02T15:04:05.000")
}

func (p *Payment) SetBankResponseTimeFromStr(value string) error {
	bankResponseTime, err := time.Parse("2006-01-02T15:04:05.000", value)
	if err != nil {
		return err
	}
	p.BankResponseTime = bankResponseTime
	return nil
}

func (p *Payment) GetBankResponseTimeStr() string {
	return p.BankResponseTime.Format("2006-01-02T15:04:05.000")
}
