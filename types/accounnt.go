package types

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        string    `json:"id"`
	Owner     string    `json:"owner_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
	DeletedAt time.Time `json:"deleted"`
}

func NewAccount(owner_id string, balance float64) *Account {
	return &Account{
		ID:        uuid.NewString(),
		Owner:     owner_id,
		Balance:   balance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: time.Now(),
	}
}

func (a *Account) ValidateAccount() bool {
	return a.Owner != "" && a.Balance >= 0
}
