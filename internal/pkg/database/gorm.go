package database

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// _db represents the global gorm.DB.
var _db *gorm.DB

func DB() *gorm.DB {
	return _db
}

func SetupGorm() error {
	cfg := config.Configs().SQLite
	db, err := gorm.Open("sqlite3", cfg.Database)
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

	err = migrate(db)
	if err != nil {
		return err
	}

	_db = db
	return nil
}

func migrate(db *gorm.DB) error {
	for _, m := range []interface{}{
		&model.User{}, &model.Filter{},
	} {
		rdb := db.AutoMigrate(m)
		if rdb.Error != nil {
			return rdb.Error
		}
	}
	return nil
}
