package entities

import "github.com/google/uuid"

type Merchant struct {
	ID       uuid.UUID
	Username string
	Password string
	Name     string
	Active   bool
	MaxQPS   int
}
