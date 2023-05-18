package memory

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	entity "github.com/thiagolcmelo/payment-gateway/merchant/entities"
	"golang.org/x/crypto/bcrypt"
)

// Storage is an in memory implementation of the Storage interface
type Storage struct {
	mechants  map[uuid.UUID]entity.Merchant
	usernames map[string]uuid.UUID
	sync.Mutex
}

var (
	// ErrUsernameAlreadyExists must be returned when trying to create a new user with an username already registered
	ErrUsernameAlreadyExists = errors.New("username already exists")
	// ErrUnknownMerchantID must be used when trying to read a merchant with an id that does not reference any merchant
	ErrUnknownMerchantID = errors.New("id does not match any merchant")
	// ErrUnknownMerchantUsername must be used when trying to read a merchant with an username that does not reference any merchant
	ErrUnknownMerchantUsername = errors.New("username does not match any merchant")
	// ErrInvalidPassword must be used when a password is incorrect for a given username
	ErrInvalidPassword = errors.New("invalid password")
)

// NewMemoryStorage is a factory for an in memory Storage
func NewMemoryStorage() *Storage {
	return &Storage{
		mechants:  make(map[uuid.UUID]entity.Merchant),
		usernames: make(map[string]uuid.UUID),
	}
}

// CreateMerchant stores a merchant in memory
func (s *Storage) CreateMerchant(m entity.Merchant) (uuid.UUID, error) {
	s.Lock()
	defer s.Unlock()

	err := m.Validate()
	if err != nil {
		return uuid.Nil, err
	}

	if _, ok := s.usernames[m.Username]; ok {
		return uuid.Nil, ErrUsernameAlreadyExists
	}

	// if ID is nil or not available, create a new one
	if _, ok := s.mechants[m.ID]; ok || m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	s.mechants[m.ID] = m

	// this will help optimize searches and validations
	s.usernames[m.Username] = m.ID

	return m.ID, nil
}

// ReadMerchant retrieves a merchant stored in memory
func (s *Storage) ReadMerchant(id uuid.UUID) (entity.Merchant, error) {
	s.Lock()
	defer s.Unlock()

	merchant, ok := s.mechants[id]

	if !ok {
		return entity.Merchant{}, ErrUnknownMerchantID
	}

	return merchant, nil
}

func (s *Storage) FindMerchantID(username string, password string) (uuid.UUID, error) {
	s.Lock()
	defer s.Unlock()

	id, ok := s.usernames[username]
	if !ok {
		return uuid.Nil, ErrUnknownMerchantUsername
	}

	// ideally this should never happen is all class invariants are respected
	merchant, ok := s.mechants[id]
	if !ok {
		return uuid.Nil, ErrUnknownMerchantID
	}

	err := bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(password))
	if err != nil {
		return uuid.Nil, ErrInvalidPassword
	}

	return id, nil
}
