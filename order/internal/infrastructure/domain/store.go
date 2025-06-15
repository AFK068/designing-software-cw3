package domain

import (
	"context"

	"github.com/google/uuid"
)

type Store interface {
	GetOrder(ctx context.Context, orderID uuid.UUID) (*Order, error)
	ListOrders(ctx context.Context, userID uuid.UUID) ([]*Order, error)
	CreateOrder(ctx context.Context, order *Order) error
	UpdateOrder(ctx context.Context, order *Order) error
	SaveOutbox(ctx context.Context, payload any) error
	GetUnprocessedOutbox(ctx context.Context, limit int) ([]OutboxRecord, error)
	MarkOutboxProcessed(ctx context.Context, id uuid.UUID) error
}
