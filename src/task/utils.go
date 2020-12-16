package task

import (
	"github.com/Aoi-hosizora/ahlib/xzone"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"sync"
	"time"
)

func checkSilent(user *model.User) bool {
	if user.Silent {
		hm, _ := xzone.MoveToZone(time.Now(), user.TimeZone)
		ss := user.SilentStart
		se := user.SilentEnd
		hour := hm.Hour()
		if ss < se { // 2 5
			if hour >= ss && hour <= se {
				return true
			}
		} else { // 22 2
			if (hour >= ss && hour <= 23) || (hour >= 0 && hour <= se) {
				return true
			}
		}
	}
	return false
}

func foreachUsers(users []*model.User, fn func(user *model.User)) {
	wg := sync.WaitGroup{}
	for _, user := range users {
		wg.Add(1)
		go func(user *model.User) {
			fn(user)
			wg.Done()
		}(user)
	}
	wg.Wait()
}
