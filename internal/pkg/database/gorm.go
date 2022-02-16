package database

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

// _db represents the global gorm.DB.
var _db *gorm.DB

func GormDB() *gorm.DB {
	return _db
}

func SetupGormDB() error {
	// open
	cfg := config.Configs().SQLite
	db, err := gorm.Open(xgorm.SQLite, xgorm.SQLiteDefaultDsn(cfg.Database))
	if err != nil {
		return err
	}

	// configure
	if !cfg.LogMode {
		db.SetLogger(xgorm.NewSilenceLogger())
	} else {
		db.LogMode(true)
		db.SetLogger(xgorm.NewLogrusLogger(logger.Logger(), xgorm.WithSlowThreshold(time.Millisecond*500)))
	}
	db.SingularTable(true)
	gorm.DefaultTableNameHandler = func(_ *gorm.DB, name string) string { return "tbl_" + name }
	xgorm.HookDeletedAt(db, xgorm.DefaultDeletedAtTimestamp)

	// migrate
	err = migrateDB(db)
	if err != nil {
		_ = db.Close()
		return err
	}

	_db = db
	return nil
}

func migrateDB(db *gorm.DB) error {
	for _, m := range []interface{}{
		&model.Chat{},
	} {
		rdb := db.AutoMigrate(m)
		if rdb.Error != nil {
			return rdb.Error
		}
	}
	return nil
}
