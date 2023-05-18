package entity

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrInvalidNumber must be used if card number is invalid
	ErrInvalidNumber = errors.New("invalid number")
	// ErrInvalidName must be used if card name is invalid
	ErrInvalidName = errors.New("invalid name")
	// ErrCardExpired must be used if card is expired
	ErrCardExpired = errors.New("invalid expiration")
	// ErrInvalidCVV must be used if card has invalid CVV
	ErrInvalidCVV = errors.New("invalid cvv")
)

type CreditCard struct {
	Number      string
	Name        string
	ExpireMonth int
	ExpireYear  int
	CVV         int
}

func NewCreditCard(number, name string, expireMonth, expireYear, cvv int) (CreditCard, error) {
	c := CreditCard{
		Number:      number,
		Name:        name,
		ExpireMonth: expireMonth,
		ExpireYear:  expireYear,
		CVV:         cvv,
	}
	return c, c.Validate()
}

func (c CreditCard) Validate() error {
	if c.Name == "" {
		return ErrInvalidName
	}
	// it must be improved
	if c.Number == "" {
		return ErrInvalidNumber
	}
	date, err := time.Parse("2006-01", fmt.Sprintf("%d-%02d", c.ExpireYear, c.ExpireMonth))
	if err != nil {
		return err
	}
	if date.Before(time.Now()) {
		return ErrCardExpired
	}
	if c.CVV < 0 || c.CVV > 999 {
		return ErrInvalidCVV
	}
	return nil
}
