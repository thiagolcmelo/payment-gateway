package entity

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUsernameOrPasswordEmpty must be used when username or password is empty
	ErrUsernameOrPasswordEmpty = errors.New("username or password empty")
	// ErrNameEmpty must be used when name is empty
	ErrNameEmpty = errors.New("name empty")
	// ErrMaxQPSCannotBeNegative must be used when MaxQPS is negative
	ErrMaxQPSCannotBeNegative = errors.New("max qps cannot be negative")
)

// Merchant represent a merchant for authentication and authorization purposes
type Merchant struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Active   bool      `json:"active"`
	MaxQPS   int       `json:"max_qps"`
}

// Validate asserts Username and Password are not empty, and that MaxQPS is not negative
func (m Merchant) Validate() error {
	if m.MaxQPS < 0 {
		return ErrMaxQPSCannotBeNegative
	}
	if m.Username == "" || m.Password == "" {
		return ErrUsernameOrPasswordEmpty
	}
	if m.Name == "" {
		return ErrNameEmpty
	}
	return nil
}

// NewMerchant is a factory for Merchant
func NewMerchant(username, password, name string, active bool, maxQPS int) (Merchant, error) {
	merchant := Merchant{
		Username: username,
		Password: password,
		Name:     name,
		Active:   active,
		MaxQPS:   maxQPS,
	}
	err := merchant.Validate()
	if err != nil {
		merchant.Password = "******"
		return merchant, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Merchant{}, nil
	}
	merchant.Password = string(hash)

	return merchant, nil
}
