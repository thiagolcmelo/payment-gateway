package storage

import (
	"github.com/google/uuid"

	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
)

// Storage defines a common interface for persisting payments
type Storage interface {
	Create(entity.Payment) (uuid.UUID, error)
	Read(uuid.UUID) (entity.Payment, error)
	ReadUsingBankReference(uuid.UUID) (entity.Payment, error)
	Update(entity.Payment) error
}
