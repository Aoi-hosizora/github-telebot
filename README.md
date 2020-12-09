# github-telebot

+ A github telegram bot built by [tucnak/telebot.v2](https://github.com/tucnak/telebot/tree/v2).

### Function

+ Query github's activity events and issue events
+ Run task with cron (through mysql and redis)
+ Send messages silently

### Endpoints

```text
*Commands*
/start - show start message
/help - show help message
/cancel - cancel the last action

*Account*
/bind - bind with a new github account
/unbind - unbind an old github account
/me - show the bound user's information
/enablesilent - enable client send
/disablesilent - disable client send

*Events*
/allowissue - allow bot to send issue
/disallowissue - allow bot to send issue
/activity - show the first page of activity events
/activityn - show the nth page of activity events
/issue - show the first page of issue events
/issuen - show the nth page of issue events
```

### References

+ [Aoi-hosizora/telebot-scaffold](https://github.com/Aoi-hosizora/telebot-scaffold)
