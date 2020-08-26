package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xredis"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var (
	Conn    redis.Conn
	redisMu sync.Mutex
)

func SetupRedis() error {
	cfg := config.Configs.Redis
	conn, err := redis.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		redis.DialPassword(cfg.Password),
		redis.DialDatabase(int(cfg.Db)),
		redis.DialConnectTimeout(time.Duration(cfg.ConnectTimeout)*time.Millisecond),
		redis.DialReadTimeout(time.Duration(cfg.ReadTimeout)*time.Millisecond),
		redis.DialWriteTimeout(time.Duration(cfg.WriteTimeout)*time.Millisecond),
	)
	if err != nil {
		return err
	}

	// Conn = xredis.NewRedisLogrus(conn, logger.Logger, config.Configs.Meta.RunMode == "debug")
	Conn = xredis.NewLogrusLogger(conn, logger.Logger, false)

	return nil
}
