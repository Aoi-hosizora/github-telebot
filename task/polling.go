package task

// var (
// 	oldStr  = ""
// 	oldObjs = make([]*model.GithubEvent, 0)
// 	dataCh  = make(chan string)
// )

func Polling() {
	// chat, err := server.Bot.ChatByID(server.Config.TelegramConfig.ChannelId)
	// if err != nil {
	// 	log.Fatalf("Failed to get chat %s:%v\n", server.Config.TelegramConfig.ChannelId, err)
	// }
	// log.Printf("Success to find channel \"%s\"\n", server.Config.TelegramConfig.ChannelId)
	//
	// for {
	// 	polling(server, chat)
	// }
}

// func polling(bot *telebot.Bot, chat *telebot.Chat) {
// 	go func() {
// 		// content, err := service.GetActions(server.Config.GithubConfig, 1)
// 		// if err != nil {
// 		dataCh <- ""
// 		// 	return
// 		// }
// 		// dataCh <- content
// 	}()
// 	newStr := <-dataCh
//
// 	if newStr != "" && newStr != oldStr {
// 		newObjs := make([]*model.GithubEvent, 0)
// 		err := json.Unmarshal([]byte(newStr), &newObjs)
// 		if err == nil {
// 			diffItf := xslice.SliceDiff(xslice.Sti(newObjs), xslice.Sti(oldObjs))
// 			diff := xslice.Its(diffItf, &model.GithubEvent{}).([]*model.GithubEvent)
// 			if len(diff) != 0 { // new
// 				msg, err := server.Bot.Send(chat, service.WrapGithubActions(diff), telebot.ModeMarkdown)
// 				if err != nil {
// 					log.Println("Failed to send message:", err)
// 				} else {
// 					log.Println("Send message success:", msg.ID)
// 				}
// 			}
// 			oldStr = newStr
// 			oldObjs = newObjs
// 		}
// 	}
// 	time.Sleep(time.Second * time.Duration(server.Config.ServerConfig.PollingDuration))
// }
