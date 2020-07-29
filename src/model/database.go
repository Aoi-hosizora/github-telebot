package model

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func SetupGorm() error {
	cfg := config.Configs.Mysql
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open("mysql", source)
	if err != nil {
		return err
	}

	db.LogMode(config.Configs.Meta.RunMode == "debug")
	db.SingularTable(true)
	db.SetLogger(xgorm.NewGormLogrus(logger.Logger))
	gorm.DefaultTableNameHandler = func(db *gorm.DB, name string) string {
		return "tbl_" + name
	}
	xgorm.HookDeleteAtField(db, xgorm.DefaultDeleteAtTimeStamp)

	err = migrate(db)
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func migrate(db *gorm.DB) error {
	for _, model := range []interface{}{
		&User{},
	} {
		rdb := db.AutoMigrate(model)
		if rdb.Error != nil {
			return rdb.Error
		}
	}
	return nil
}
