package infrastructure

import (
	"context"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/higansama/xyz-multi-finance/config"
	"github.com/higansama/xyz-multi-finance/internal/app"
	"github.com/higansama/xyz-multi-finance/internal/cache"
	"github.com/higansama/xyz-multi-finance/internal/db"
	"github.com/higansama/xyz-multi-finance/internal/events"
	"github.com/higansama/xyz-multi-finance/internal/inbox"
	"github.com/higansama/xyz-multi-finance/internal/locker"
	"github.com/higansama/xyz-multi-finance/internal/outbox"
	"github.com/higansama/xyz-multi-finance/internal/processor"
	"github.com/higansama/xyz-multi-finance/internal/redis"
	"github.com/higansama/xyz-multi-finance/internal/request"

	// "github.com/higansama/xyz-multi-finance/internal/template"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/higansama/xyz-multi-finance/internal/validator"
	"github.com/pkg/errors"
	amqpgo "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Infrastructure struct {
	Env             app.Environment
	Config          config.Config
	MongoDB         *mongo.Database
	Redis           *redis.Client
	Cache           *cache.Cache
	Locker          *locker.Locker
	Amqp            *amqpgo.Connection
	TxManager       db.TransactionManager
	Validator       *validator.Validator
	Conforms        *mold.Transformer
	StructProcessor processor.StructProcessor
	Request         request.Request
	Middleware      Middleware
	InboxRepo       inbox.Repository
	OutboxRepo      outbox.Repository
	EventBus        *events.EventBus
	Ctx             context.Context
	ErrorCh         chan error
}

func NewInfrastructure(cfg config.Config) *Infrastructure {
	return &Infrastructure{Config: cfg}
}

func (infra *Infrastructure) InitInfrastructure(ctx context.Context) (*Infrastructure, error, func()) {
	var cleanup []func()

	infra.Ctx = ctx
	infra.ErrorCh = make(chan error)

	if infra.Config.App.Name == "" {
		return nil, errors.New("app name is required"), nil
	}
	infra.Config.App.Name = utils.KebabCase(infra.Config.App.Name)

	env, err := app.NewEnvironmentFromString(infra.Config.App.Env)
	if err != nil {
		return nil, err, nil
	}
	infra.Env = env

	mongoDB, err, mongoCleanup := infra.configMongoDB()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, mongoCleanup)
	infra.MongoDB = mongoDB

	redisClient, err, redisCleanup := infra.configRedis()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, redisCleanup)
	infra.Redis = redisClient

	redisStore := cache.NewRedisStore(redisClient)
	infra.Cache = cache.NewCache(redisStore)
	infra.Locker = locker.NewLocker(locker.NewRedisLocker(redisClient))

	amqpConn, err, amqpCleanup := infra.configAmqp()
	if err != nil {
		return nil, err, nil
	}
	cleanup = append(cleanup, amqpCleanup)
	infra.Amqp = amqpConn

	infra.Conforms = modifiers.New()

	infra.Validator = validator.ConfigValidator()

	infra.StructProcessor = processor.NewStructProcessor(infra.Conforms, infra.Validator)

	infra.Request = request.NewRequest(infra.StructProcessor)

	middleware, err := infra.setupMiddleware()
	if err != nil {
		return nil, err, nil
	}
	infra.Middleware = middleware

	infra.InboxRepo = inbox.NewMongoRepository(infra.Env, infra.MongoDB)

	infra.OutboxRepo = outbox.NewMongoRepository(infra.Env, infra.MongoDB)

	infra.TxManager = db.NewMongoTransactionManager(infra.MongoDB)

	eventBus, err := infra.setupEventBus()
	if err != nil {
		return nil, err, nil
	}
	infra.EventBus = eventBus

	// err = template.LoadTemplates()
	// if err != nil {
	// 	return nil, err, nil
	// }

	return infra, nil, func() {
		log.Info().Msg("Infra cleanup")
		for _, c := range cleanup {
			c()
		}
	}
}
