package db

import (
	"context"
	"time"

	"github.com/higansama/xyz-multi-finance/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
)

type transactionKey struct {
}

type AlwaysRollbackKey struct {
}

var transactionTimeout = 120 * time.Second

type TransactionManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
	WithTemporaryTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}

type mongoTransactionManager struct {
	db *mongo.Database
}

func NewMongoTransactionManager(mdb *mongo.Database) TransactionManager {
	return &mongoTransactionManager{mdb}
}

func (tm *mongoTransactionManager) WithTransaction(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) error {
	if _, ok := ctx.Value(transactionKey{}).(bool); ok {
		// if already in transaction, we just need to pass original ctx, and we're done.
		err := fn(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	session, err := tm.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	timeout := time.NewTimer(transactionTimeout)
	defer timeout.Stop()

	for {
		err = session.StartTransaction()
		if err != nil {
			return err
		}

		sessCtx := mongo.NewSessionContext(ctx, session)
		nCtx := context.WithValue(sessCtx, transactionKey{}, true)

		err = fn(nCtx)

		if _, ok := ctx.Value(AlwaysRollbackKey{}).(bool); ok {
			// useful in tests
			_ = session.AbortTransaction(mongodb.NewBackgroundContext(ctx))
			return err
		}

		if err != nil {
			_ = session.AbortTransaction(mongodb.NewBackgroundContext(ctx))

			select {
			case <-timeout.C:
				return err
			default:
			}

			if cerr, ok := err.(mongo.CommandError); ok {
				if cerr.HasErrorLabel(driver.TransientTransactionError) {
					continue
				}
			}

			return err
		}

		if ctx.Err() != nil {
			_ = session.AbortTransaction(mongodb.NewBackgroundContext(ctx))
			return ctx.Err()
		}

	CommitLoop:
		for {
			err = session.CommitTransaction(ctx)
			if err == nil {
				return nil
			}

			select {
			case <-timeout.C:
				return err
			default:
			}

			if cerr, ok := err.(mongo.CommandError); ok {
				if cerr.HasErrorLabel(driver.UnknownTransactionCommitResult) && !cerr.IsMaxTimeMSExpiredError() {
					continue
				}
				if cerr.HasErrorLabel(driver.TransientTransactionError) {
					break CommitLoop
				}
			}
			return err
		}
	}
}

func (tm *mongoTransactionManager) WithTemporaryTransaction(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) error {
	ctx = context.WithValue(ctx, AlwaysRollbackKey{}, true)

	err := tm.WithTransaction(ctx, fn)
	if err != nil {
		return err
	}

	return nil
}
