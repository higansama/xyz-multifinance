package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/higansama/xyz-multi-finance/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 30 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

type MongoDB struct {
	MongoClient *mongo.Client
}

// NewMongoDB Create new MongoDB client
func NewMongoDB(ctx context.Context, cfg config.MongoDB) (*MongoDB, error) {
	opt := options.Client().
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	if cfg.URI != "" {
		opt = opt.ApplyURI(cfg.URI)
	} else {
		opt = opt.SetAuth(options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: cfg.Db,
		}).SetHosts([]string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)})
	}

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.WithStack(err)
	}

	return &MongoDB{MongoClient: client}, nil
}

func (m *MongoDB) Close() error {
	return m.MongoClient.Disconnect(nil)
}

func (infra *Infrastructure) configMongoDB() (*mongo.Database, error, func()) {
	if infra.Config.DB.MongoDB.Db == "" {
		return nil, errors.New("mongo database must be set"), nil
	}

	mdb, err := NewMongoDB(infra.Ctx, infra.Config.DB.MongoDB)
	if err != nil {
		return nil, err, nil
	}

	log.Info().Msg("Mongo connected")
	db := mdb.MongoClient.Database(infra.Config.DB.MongoDB.Db)

	return db, nil, func() {
		_ = mdb.Close()
	}
}
