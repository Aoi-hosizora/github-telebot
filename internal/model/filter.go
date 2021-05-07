package model

import (
	"github.com/Aoi-hosizora/ahlib-db/xgorm"
	"time"
)

type Filter struct {
	Fid       uint64 `gorm:"primary_key; auto_increment"`
	ChatID    int64  `gorm:"type:bigint(20);   not_null; unique_index:uk_chat_user_repo_event"`
	Username  string `gorm:"type:varchar(255); not_null; unique_index:uk_chat_user_repo_event"`
	RepoName  string `gorm:"type:varchar(255); not_null; unique_index:uk_chat_user_repo_event"`
	EventType string `gorm:"type:varchar(255); not_null; unique_index:uk_chat_user_repo_event"`

	xgorm.GormTime2
	DeletedAt *time.Time `gorm:"default:'1970-01-01 00:00:01'; unique_index:uk_chat_user_repo_event"`
}

/*
CREATE TABLE "tbl_filter" (
    "fid"        integer primary key autoincrement,
    "chat_id"    bigint(20),
    "username"   varchar(255),
    "repo_name"  varchar(255),
    "event_type" varchar(255),
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime DEFAULT '1970-01-01 00:00:01'
);
CREATE UNIQUE INDEX uk_chat_user_repo_event ON "tbl_filter" (chat_id, username, repo_name, event_type, deleted_at);
*/
