package model

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"time"
)

type User struct {
	Id          uint32 `gorm:"primary_key; auto_increment"`
	ChatID      int64  `gorm:"type:bigint(20);   not_null; unique_index:uk_chat_id"`
	Username    string `gorm:"type:varchar(255); not_null"`
	Token       string `gorm:"type:varchar(255); not_null"`
	AllowIssue  bool   `gorm:"type:tinyint;      not_null; default:0"`
	FilterMe    bool   `gorm:"type:tinyint;      not_null; default:1"`
	Silent      bool   `gorm:"type:tinyint;      non_null; default:0"`
	SilentStart int    `gorm:"type:tinyint;      non_null; default:0"`
	SilentEnd   int    `gorm:"type:tinyint;      non_null; default:0"`
	TimeZone    string `gorm:"type:varchar(15);  non_null; default:'+00:00'"`

	xgorm.GormTime2
	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:01'; unique_index:uk_chat_id"`
}

/*
CREATE TABLE IF NOT EXISTS "tbl_user" (
    "id"           integer primary key autoincrement,
    "chat_id"      bigint(20),
    "username"     varchar(255),
    "token"        varchar(255),
    "allow_issue"  tinyint     DEFAULT 0,
    "filter_me"    tinyint     DEFAULT 1,
    "silent"       tinyint     DEFAULT 0,
    "silent_start" tinyint     DEFAULT 0,
    "silent_end"   tinyint     DEFAULT 0,
    "time_zone"    varchar(15) DEFAULT '+00:00',
    "created_at"   datetime,
    "updated_at"   datetime,
    "deleted_at"   datetime DEFAULT '1970-01-01 00:00:01'
);
CREATE UNIQUE INDEX uk_chat_id ON "tbl_user" (chat_id, deleted_at);
*/
