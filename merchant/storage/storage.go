package storage

import (
	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/merchant/entities"
)

// Storage defines a common interface for persisting or in memory storage for storing merchants
type Storage interface {
	CreateMerchant(entity.Merchant) (uuid.UUID, error)
	ReadMerchant(uuid.UUID) (entity.Merchant, error)
	FindMerchantID(string, string) (uuid.UUID, error)
}
