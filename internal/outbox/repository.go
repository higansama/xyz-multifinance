package outbox

import (
	"context"

	"github.com/higansama/xyz-multi-finance/internal/app"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Save(ctx context.Context, outbox *Table) error
}

type mongoRepository struct {
	env app.Environment
	db  *mongo.Database
}

func NewMongoRepository(env app.Environment, db *mongo.Database) Repository {
	// silently create collection
	_ = db.CreateCollection(context.Background(), utils.FormatNameForEnv(env, "outbox"))

	return &mongoRepository{env, db}
}

func (r *mongoRepository) Save(
	ctx context.Context,
	outbox *Table,
) error {
	col := r.db.Collection(utils.FormatNameForEnv(r.env, "outbox"))
	res, err := col.InsertOne(ctx, outbox)
	if err != nil {
		return errors.WithStack(err)
	}
	outbox.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return nil
}
