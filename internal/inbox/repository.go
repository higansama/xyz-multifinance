package inbox

import (
	"context"

	"github.com/higansama/xyz-multi-finance/internal/app"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Save(ctx context.Context, inbox *Inbox) error
	FindByMessageIdAndConsumer(ctx context.Context, msgId string, consumer string) (*Inbox, error)
}

type mongoRepository struct {
	env app.Environment
	db  *mongo.Database
}

func NewMongoRepository(env app.Environment, db *mongo.Database) Repository {
	colName := utils.FormatNameForEnv(env, "inbox")
	// silently create collection
	_ = db.CreateCollection(context.Background(), colName)
	col := db.Collection(colName)
	// silently create indexes
	mod := mongo.IndexModel{
		Keys: bson.D{
			{Key: "message_id", Value: 1},
			{Key: "consumer", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, _ = col.Indexes().CreateOne(context.Background(), mod)

	return &mongoRepository{env, db}
}

func (r *mongoRepository) Save(
	ctx context.Context,
	inbox *Inbox,
) error {
	col := r.db.Collection(utils.FormatNameForEnv(r.env, "inbox"))

	res, err := col.InsertOne(ctx, inbox)
	if err != nil {
		return errors.WithStack(err)
	}
	inbox.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return nil
}

func (r *mongoRepository) FindByMessageIdAndConsumer(
	ctx context.Context,
	msgId string,
	consumer string,
) (*Inbox, error) {
	col := r.db.Collection(utils.FormatNameForEnv(r.env, "inbox"))
	res := col.FindOne(ctx, bson.D{
		{Key: "message_id", Value: msgId},
		{Key: "consumer", Value: consumer},
	})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, errors.WithStack(res.Err())
	}

	var inbox Inbox
	err := res.Decode(&inbox)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &inbox, nil
}
