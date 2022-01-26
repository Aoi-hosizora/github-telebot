module github.com/Aoi-hosizora/github-telebot

go 1.14

require (
	github.com/Aoi-hosizora/ahlib v0.0.0-00010101000000-000000000000
	github.com/Aoi-hosizora/ahlib-db/xgorm v0.0.0-00010101000000-000000000000
	github.com/Aoi-hosizora/ahlib-db/xredis v0.0.0-00010101000000-000000000000
	github.com/Aoi-hosizora/ahlib-more v0.0.0-00010101000000-000000000000
	github.com/Aoi-hosizora/ahlib-web v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.4.11
	github.com/jinzhu/gorm v1.9.16
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/tucnak/telebot.v2 v2.5.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/Aoi-hosizora/ahlib => ../ahlib
	github.com/Aoi-hosizora/ahlib-db/xgorm => ../ahlib-db/xgorm
	github.com/Aoi-hosizora/ahlib-db/xredis => ../ahlib-db/xredis
	github.com/Aoi-hosizora/ahlib-more => ../ahlib-more
	github.com/Aoi-hosizora/ahlib-web => ../ahlib-web
)
