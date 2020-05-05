package model

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"github.com/Aoi-hosizora/ahlib-web/xgorm"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func SetupGorm() error {
	cfg := config.Configs.MysqlConfig
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := gorm.Open("mysql", source)
	if err != nil {
		return err
	}

	db.LogMode(gin.Mode() == gin.DebugMode)
	db.SingularTable(true)
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
	rdb := db.AutoMigrate(&User{})
	if rdb.Error != nil {
		return rdb.Error
	}
	return nil
}
