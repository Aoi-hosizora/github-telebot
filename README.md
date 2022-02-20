# github-telebot

+ A GitHub event notifier telegram bot built by [tucnak/telebot.v2](https://github.com/tucnak/telebot/tree/v2).

### Function

+ [x] Notify new activity and issue events for subscribed GitHub account
+ [x] Configurable Silent and NoPreview send options 

### Endpoints

```text
*Start*
/start - show start message
/help - show this help message
/cancel - cancel the last action

*Subscribe*
/subscribe - subscribe with a new GitHub account
/unsubscribe - unsubscribe the current GitHub account
/me - show the subscribed user's information

*Option*
/allowissue - allow bot to notify new issue events
/disallowissue - disallow bot to notify new issue events
/enablesilent - send message with no notification
/disablesilent - send message with notification
/enablepreview - enable preview for link
/disablepreview - disable preview for link

*Event*
/activity - show the first page of activity events
/activity N - show the N-th page of activity events
/issue - show the first page of issue events
/issue N - show the N-th page of issue events
```
