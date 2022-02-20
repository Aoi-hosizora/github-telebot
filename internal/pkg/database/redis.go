package database

import (
	"context"
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xredis"
	"github.com/Aoi-hosizora/ahlib/xpointer"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/go-redis/redis/v8"
	"runtime"
	"time"
)

var _redis *redis.Client

func RedisClient() *redis.Client {
	return _redis
}

func SetupRedisClient() error {
	// open
	cfg := config.Configs().Redis
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:       int(cfg.DB),
		Password: cfg.Password,

		DialTimeout:  time.Duration(xpointer.Int32Val(cfg.DialTimeout, 5)) * time.Second,  // defaults to 5s (cannot set to unlimited)
		ReadTimeout:  time.Duration(xpointer.Int32Val(cfg.ReadTimeout, 3)) * time.Second,  // defaults to 3s (cannot set to unlimited)
		WriteTimeout: time.Duration(xpointer.Int32Val(cfg.WriteTimeout, 3)) * time.Second, // defaults to 3s (cannot set to unlimited)
		OnConnect:    nil,

		PoolSize:     int(xpointer.Int32Val(cfg.MaxOpens, int32(10*runtime.NumCPU()))),    // defaults to 10 * #CPU (cannot set to unlimited)
		MinIdleConns: int(xpointer.Int32Val(cfg.MinIdles, 1)),                             // defaults to 1 (can set to 0)
		MaxConnAge:   time.Duration(xpointer.Int32Val(cfg.MaxLifetime, 0)) * time.Second,  // defaults to unlimited
		IdleTimeout:  time.Duration(xpointer.Int32Val(cfg.MaxIdletime, -1)) * time.Second, // defaults to unlimited
	}
	client := redis.NewClient(opt)

	// test ping
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return err
	}

	// configure
	if !cfg.LogMode {
		redis.SetLogger(xredis.NewSilenceLogger())
	} else {
		redis.SetLogger(xredis.NewSilenceLogger())
		client.AddHook(xredis.NewLogrusLogger(logger.Logger(), xredis.WithSlowThreshold(time.Millisecond*100)))
	}

	_redis = client
	return nil
}
