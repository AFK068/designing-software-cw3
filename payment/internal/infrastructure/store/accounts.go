package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/AFK068/designing-software-cw3/payment/internal/domain"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/store/txs"
)

func (s *Store) CreateAccount(ctx context.Context, account *domain.Account) error {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `INSERT INTO accounts (user_id, amount) VALUES ($1, $2)`

	if _, err := querier.Exec(ctx, query, account.UserID, account.Amount); err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (s *Store) UpdateAccount(ctx context.Context, account *domain.Account) error {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `UPDATE accounts SET amount = $1 WHERE user_id = $2`

	if _, err := querier.Exec(ctx, query, account.Amount, account.UserID); err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (s *Store) GetAccount(ctx context.Context, userID uuid.UUID) (*domain.Account, error) {
	querier := txs.GetQuerier(ctx, s.pool)

	query := `SELECT amount FROM accounts WHERE user_id = $1`
	row := querier.QueryRow(ctx, query, userID)

	account := &domain.Account{}
	if err := row.Scan(&account.Amount); err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("account not found for user ID %s: %w", userID, err)
		}

		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	account.UserID = userID

	return account, nil
}
