package memory

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/ledger/entities"
)

var (
	// ErrUnknownPayment must be use while trying to read payment with unknown id
	ErrUnknownPayment = errors.New("there is no payment with given id")
	// ErrUnknownBankReference must be used when trying to read using unknown bank reference
	ErrUnknownBankReference = errors.New("unknown bank reference")
)

// Storage is an in memory implementation of a Ledger to store payments
type Storage struct {
	payments       map[uuid.UUID]entity.Payment
	bankReferences map[uuid.UUID]uuid.UUID
	sync.RWMutex
}

// NewMemoryStorage is a factory for in memory Storage for payments
func NewMemoryStorage() *Storage {
	return &Storage{
		payments:       make(map[uuid.UUID]entity.Payment),
		bankReferences: make(map[uuid.UUID]uuid.UUID),
	}
}

// Create adds a new payment to the Ledger
func (l *Storage) Create(p entity.Payment) (uuid.UUID, error) {
	l.Lock()
	defer l.Unlock()

	var id uuid.UUID

	for {
		id = uuid.New()
		if _, ok := l.payments[id]; !ok {
			break
		}
	}

	p.ID = id

	err := p.Validate()
	if err != nil {
		return uuid.Nil, err
	}

	l.payments[id] = p
	if p.BankPaymentID != uuid.Nil {
		l.bankReferences[p.BankPaymentID] = p.ID
	}

	return id, nil
}

// Read returns details of a payment in the Ledger
func (l *Storage) Read(id uuid.UUID) (entity.Payment, error) {
	l.RLock()
	defer l.RUnlock()

	payment, ok := l.payments[id]
	if !ok {
		return entity.Payment{}, ErrUnknownPayment
	}

	return payment, nil
}

// ReadUsingBankReference returns details of a payment in the Ledger using a bank reference
func (l *Storage) ReadUsingBankReference(id uuid.UUID) (entity.Payment, error) {
	l.RLock()
	defer l.RUnlock()

	paymentID, ok := l.bankReferences[id]
	if !ok {
		return entity.Payment{}, ErrUnknownBankReference
	}

	payment, ok := l.payments[paymentID]
	if !ok {
		return entity.Payment{}, ErrUnknownPayment
	}

	return payment, nil
}

// Update edits information of a given payment
func (l *Storage) Update(p entity.Payment) error {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.payments[p.ID]; !ok {
		return ErrUnknownPayment
	}

	err := p.Validate()
	if err != nil {
		return err
	}

	l.payments[p.ID] = p
	if p.BankPaymentID != uuid.Nil {
		l.bankReferences[p.BankPaymentID] = p.ID
	}

	return nil
}
