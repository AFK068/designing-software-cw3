package domain

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	CreateAccount(ctx context.Context, account *Account) error
	UpdateAccount(ctx context.Context, account *Account) error
	GetAccount(ctx context.Context, userID uuid.UUID) (*Account, error)
}
