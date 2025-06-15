package txs

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

type TxBeginner struct {
	db *pgxpool.Pool
}

func NewTxBeginner(db *pgxpool.Pool) *TxBeginner {
	return &TxBeginner{
		db: db,
	}
}

func (t *TxBeginner) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) (err error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback(ctx))
		}
	}()

	err = txFunc(injectTx(ctx, tx))
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
