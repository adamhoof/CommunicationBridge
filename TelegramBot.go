package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

var Bot *tb.Bot

type TelegramBotHandler struct {
	keyboards map[string]*tb.ReplyMarkup
}

type User struct {
	userId string
}

func (user *User) Recipient() string {
	return user.userId
}

func (botHandler *TelegramBotHandler) CreateBot() {
	var err error
	Bot, err = tb.NewBot(tb.Settings{
		Token: "1914152683:AAF4r5URK9fCoJicXsCADukXuiTQSYM--U8",
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}
	botHandler.keyboards = make(map[string]*tb.ReplyMarkup)
}

func CreateHumanReadable(applianceDataMap map[string]interface{}) string {

	var humanReadable string

	if applianceDataMap != nil {

		for key, value := range applianceDataMap {
			humanReadable += key + ": " + value.(string) + "\n"
		}
		return humanReadable
	}
	return "map iterating yeeted"
}

func (botHandler *TelegramBotHandler) StartBot() {
	Bot.Start()
}

func SendMessage(usr User, title string, message interface{}) {

	_, err := Bot.Send(&usr, title, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}
