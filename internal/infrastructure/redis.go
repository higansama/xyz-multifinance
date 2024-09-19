package infrastructure

import (
	"github.com/higansama/xyz-multi-finance/internal/redis"
	"github.com/higansama/xyz-multi-finance/internal/utils"
	"github.com/rs/zerolog/log"
)

func (infra *Infrastructure) configRedis() (*redis.Client, error, func()) {
	if infra.Config.DB.Redis.Prefix == "" {
		infra.Config.DB.Redis.Prefix = utils.FormatNameForEnv(infra.Env, infra.Config.App.Name)
	}

	rds, err := redis.NewRedisClient(infra.Ctx, infra.Config.DB.Redis)
	if err != nil {
		return nil, err, nil
	}

	log.Info().Msg("Redis connected")

	return rds, nil, func() {
		_ = rds.Close()
	}
}
