package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xredis"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	Rpool *redis.Pool
)

func SetupRedis() error {
	cfg := config.Configs.Redis
	Rpool = &redis.Pool{
		MaxIdle:         int(cfg.MaxIdle),
		MaxActive:       int(cfg.MaxActive),
		MaxConnLifetime: time.Duration(cfg.MaxLifetime) * time.Second,
		IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(
				"tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				redis.DialPassword(cfg.Password),
				redis.DialDatabase(int(cfg.Db)),
				redis.DialConnectTimeout(time.Duration(cfg.ConnectTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
			)
			if err != nil {
				return nil, err
			}

			conn = xredis.NewLogrusRedis(conn, logger.Logger, false).WithSkip(4)
			conn = xredis.NewMutexRedis(conn)
			return conn, nil
		},
	}

	conn := Rpool.Get()
	defer conn.Close()
	err := conn.Err()
	if err != nil {
		return err
	}

	return nil
}
