package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type TelegramBotHandler struct {
	bot *tb.Bot
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
	botHandler.bot, err = tb.NewBot(tb.Settings{
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
	botHandler.bot.Start()
}

func SendMessage(telegramBotHandler *TelegramBotHandler, usr User, title string, message interface{}) {

	_, err := telegramBotHandler.bot.Send(&usr, title, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}
