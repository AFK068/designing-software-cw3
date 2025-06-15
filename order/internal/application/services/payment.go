package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"
)

type Transactor interface {
	WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type OrderService struct {
	store      domain.Store
	transactor Transactor
}

func NewOrderService(store domain.Store, transactor Transactor) *OrderService {
	return &OrderService{
		store:      store,
		transactor: transactor,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID, amount int64, description string) (*domain.Order, error) {
	var order *domain.Order

	err := s.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		order = &domain.Order{
			ID:          uuid.New(),
			UserID:      userID,
			Amount:      amount,
			Description: description,
			Status:      domain.StatusNew,
		}
		if err := s.store.CreateOrder(ctx, order); err != nil {
			return err
		}

		outboxPayload := map[string]any{
			"event":       "order_created",
			"order_id":    order.ID,
			"user_id":     order.UserID,
			"amount":      order.Amount,
			"description": order.Description,
		}
		if err := s.store.SaveOutbox(ctx, outboxPayload); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) ListOrders(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	return s.store.ListOrders(ctx, userID)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	return s.store.GetOrder(ctx, orderID)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return s.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		order, err := s.store.GetOrder(ctx, orderID)
		if err != nil {
			return err
		}

		order.Status = status

		return s.store.UpdateOrder(ctx, order)
	})
}
