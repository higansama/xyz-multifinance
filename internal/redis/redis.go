package redis

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
	"strings"

	"github.com/higansama/xyz-multi-finance/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Prefix   string
	Instance *redis.Client
}

func NewRedisClient(ctx context.Context, cfg config.Redis) (*Client, error) {
	db, err := strconv.Atoi(cfg.Db)
	if err != nil {
		return nil, errors.Errorf("redis: invalid database number: %q", db)
	}

	opts := &redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       db,
	}
	if cfg.Secure {
		opts.TLSConfig = &tls.Config{
			ServerName: cfg.Host,
		}
	}

	if cfg.URI != "" {
		opt, err := redis.ParseURL(cfg.URI)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		opts = opt
	}

	if !strings.HasSuffix(cfg.Prefix, ":") {
		cfg.Prefix += ":"
	}

	instance := redis.NewClient(opts)
	if res := instance.Ping(ctx); res.Err() != nil {
		return nil, errors.Wrap(res.Err(), "redis: could not establish connection")
	}

	return &Client{Instance: instance, Prefix: cfg.Prefix}, nil
}

func (c *Client) Close() error {
	return c.Instance.Close()
}
