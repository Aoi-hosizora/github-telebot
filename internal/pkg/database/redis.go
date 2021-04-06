package database

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xredis"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/go-redis/redis/v8"
	"time"
)

// _redis represents the global redis.Client.
var _redis *redis.Client

func Redis() *redis.Client {
	return _redis
}

func SetupRedis() error {
	cfg := config.Configs().Redis
	client := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:           int(cfg.DB),
		Password:     cfg.Password,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,  // defaults to 5s
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,  // defaults to 3s
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second, // defaults to 3s

		PoolSize:    int(cfg.MaxOpen),                             // defaults to 10 * #CPU
		MaxConnAge:  time.Duration(cfg.MaxLifetime) * time.Second, // defaults to unlimited
		IdleTimeout: time.Duration(cfg.MaxIdletime) * time.Second, // defaults to 5min
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	redis.SetLogger(xredis.NewSilenceLogger())
	if cfg.LogMode {
		client.AddHook(xredis.NewLogrusLogger(logger.Logger()))
	}

	_redis = client
	return nil
}
