package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

func DbTransaction(
	db *mongo.Database,
	ctx context.Context,
	fn func(ctx mongo.SessionContext) error,
) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (any, error) {
		err = fn(sessCtx)
		if err != nil {
			return nil, fmt.Errorf("transaction@DbTransaction->fn: %w", err)
		}

		return nil, nil
	})

	return err
}
