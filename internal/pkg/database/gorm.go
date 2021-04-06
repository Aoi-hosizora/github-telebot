package database

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

// _db represents the global gorm.DB.
var _db *gorm.DB

func DB() *gorm.DB {
	return _db
}

func SetupGorm() error {
	cfg := config.Configs().MySQL
	dsl := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)
	db, err := gorm.Open("mysql", dsl)
	if err != nil {
		return err
	}

	if cfg.LogMode {
		db.LogMode(true)
		db.SetLogger(xgorm.NewLogrusLogger(logger.Logger()))
	} else {
		db.SetLogger(xgorm.NewSilenceLogger())
	}
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "tbl_" + defaultTableName
	}
	xgorm.HookDeletedAt(db, xgorm.DefaultDeletedAtTimestamp)

	db.DB().SetMaxOpenConns(int(cfg.MaxOpen))                                // defaults to unlimited
	db.DB().SetMaxIdleConns(int(cfg.MaxIdle))                                // defaults to 2
	db.DB().SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second) // defaults to unlimited
	db.DB().SetConnMaxIdleTime(time.Duration(cfg.MaxIdletime) * time.Second) // defaults to unlimited

	err = migrate(db)
	if err != nil {
		return err
	}

	_db = db
	return nil
}

func migrate(db *gorm.DB) error {
	for _, m := range []interface{}{
		&model.User{},
	} {
		rdb := db.AutoMigrate(m)
		if rdb.Error != nil {
			return rdb.Error
		}
	}
	return nil
}
