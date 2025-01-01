package transaction

import (
	"context"

	"github.com/igortoigildin/goph-keeper/internal/server/client/db"
	"github.com/igortoigildin/goph-keeper/internal/server/client/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager creates new manager, which implements db.TxManager interface
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	// If this is intermediate tx, skip initiation of the new tx
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	// Start new tx
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction")
	}

	// Put tx into context
	ctx = pg.MakeContextTx(ctx, tx)

	// Set defer func for tx roll back or commit accordingly
	defer func() {
		// recover after panic
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		// case with not nil error
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}

			return
		}

		// case with nil error
		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}
	}()

	// Init original fn
	// In case of any error - return error, and defer func will roll back,
	// otherwise, tx will be committed.
	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside tx")
	}

	return err
}

