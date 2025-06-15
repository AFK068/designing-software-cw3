package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/store/txs"
	"github.com/google/uuid"
)

func (s *Store) SaveOutbox(ctx context.Context, payload any) error {
	querier := txs.GetQuerier(ctx, s.pool)

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal outbox payload: %w", err)
	}

	query := `INSERT INTO outbox (payload) VALUES ($1)`

	_, err = querier.Exec(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to insert outbox record: %w", err)
	}

	return nil
}

func (s *Store) GetUnprocessedOutbox(ctx context.Context, limit int) ([]domain.OutboxRecord, error) {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `SELECT id, payload FROM outbox WHERE processed = FALSE ORDER BY created_at LIMIT $1`

	rows, err := querier.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get outbox records: %w", err)
	}
	defer rows.Close()

	var result []domain.OutboxRecord

	for rows.Next() {
		var rec domain.OutboxRecord
		if err := rows.Scan(&rec.ID, &rec.Payload); err != nil {
			return nil, fmt.Errorf("failed to scan outbox record: %w", err)
		}

		result = append(result, rec)
	}

	return result, nil
}

func (s *Store) MarkOutboxProcessed(ctx context.Context, id uuid.UUID) error {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `UPDATE outbox SET processed = TRUE WHERE id = $1`

	_, err := querier.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark outbox as processed: %w", err)
	}

	return nil
}
