package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/AFK068/designing-software-cw3/payment/internal/domain"
)

type Transactor interface {
	WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error
}

type PaymentService struct {
	store      domain.Store
	transactor Transactor
}

func NewPaymentService(store domain.Store, transactor Transactor) *PaymentService {
	return &PaymentService{
		store:      store,
		transactor: transactor,
	}
}

func (s *PaymentService) GetAccount(ctx context.Context, userID uuid.UUID) (*domain.Account, error) {
	return s.store.GetAccount(ctx, userID)
}

func (s *PaymentService) CreateAccount(ctx context.Context, userID uuid.UUID) (*domain.Account, error) {
	account := &domain.Account{
		UserID: userID,
		Amount: 0,
	}

	err := s.store.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PaymentService) ReplenishAccount(ctx context.Context, userID uuid.UUID, amount int64) (*domain.Account, error) {
	var result *domain.Account

	err := s.transactor.WithTransaction(ctx, func(ctx context.Context) error {
		account, err := s.store.GetAccount(ctx, userID)
		if err != nil {
			return err
		}

		account.Amount += amount

		if err := s.store.UpdateAccount(ctx, account); err != nil {
			return err
		}

		result = account

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
