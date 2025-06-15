package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/store/txs"
)

func (s *Store) CreateOrder(ctx context.Context, order *domain.Order) error {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `INSERT INTO orders (id, user_id, description, amount, status) VALUES ($1, $2, $3, $4, $5)`

	_, err := querier.Exec(ctx, query, order.ID, order.UserID, order.Description, order.Amount, order.Status)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (s *Store) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `SELECT id, user_id, description, amount, status FROM orders WHERE id = $1`

	row := querier.QueryRow(ctx, query, orderID)

	var order domain.Order

	var status string
	if err := row.Scan(&order.ID, &order.UserID, &order.Description, &order.Amount, &status); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %w", err)
		}

		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.Status = domain.OrderStatus(status)

	return &order, nil
}

func (s *Store) ListOrders(ctx context.Context, userID uuid.UUID) ([]*domain.Order, error) {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `SELECT id, user_id, description, amount, status FROM orders WHERE user_id = $1`

	rows, err := querier.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {
		var order domain.Order

		var status string

		if err := rows.Scan(&order.ID, &order.UserID, &order.Description, &order.Amount, &status); err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		order.Status = domain.OrderStatus(status)
		orders = append(orders, &order)
	}

	return orders, nil
}

func (s *Store) UpdateOrder(ctx context.Context, order *domain.Order) error {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `UPDATE orders SET description = $1, amount = $2, status = $3 WHERE id = $4`

	_, err := querier.Exec(ctx, query, order.Description, order.Amount, order.Status, order.ID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}
