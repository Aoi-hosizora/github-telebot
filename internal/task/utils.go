package task

import (
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/Aoi-hosizora/github-telebot/internal/model"
	"sync"
	"time"
)

func checkSilent(user *model.User) bool {
	if !user.Silent {
		return false
	}

	loc, err := xtime.ParseTimezone(user.TimeZone)
	if err != nil {
		return false
	}
	hm := time.Now().In(loc)

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
