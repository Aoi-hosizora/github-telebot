package model

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"time"
)

type Chat struct {
	Cid      uint64 `gorm:"primary_key; auto_increment"`
	ChatID   int64  `gorm:"type:bigint(20);   not_null; unique_index:uk_chat_id"`
	Username string `gorm:"type:varchar(255); not_null"`
	Token    string `gorm:"type:varchar(255); not_null"`
	Issue    bool   `gorm:"type:tinyint(1);   not_null; default:0"`
	FilterMe bool   `gorm:"type:tinyint(1);   not_null; default:1"`
	Silent   bool   `gorm:"type:tinyint(1);   non_null; default:0"`
	Preview  bool   `gorm:"type:tinyint(1);   not_null; default:1"`

	xgorm.GormTime2
	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:01'; unique_index:uk_chat_id"`
}

/*
CREATE TABLE "tbl_chat" (
    "cid"        integer primary key autoincrement,
    "chat_id"    bigint(20),
    "username"   varchar(255),
    "issue"      tinyint(1) DEFAULT 0,
    "filter_me"  tinyint(1) DEFAULT 1,
    "silent"     tinyint(1) DEFAULT 0,
    "preview"    tinyint(1) DEFAULT 1,
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime DEFAULT '1970-01-01 00:00:01'
);
CREATE UNIQUE INDEX uk_chat_id ON "tbl_chat" (chat_id, deleted_at);
*/
