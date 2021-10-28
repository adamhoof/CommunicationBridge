package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"time"
)

type TelegramBot struct {
	bot       *tb.Bot
	keyboards map[string]*tb.ReplyMarkup
}

type User struct {
	id string
}

const meId = "558297691"

var me = User{id: meId}

const (
	TXT    = 0
	KBOARD = 1
)

func (user *User) Recipient() string {
	return user.id
}

func (telegramBot *TelegramBot) CreateBot() {

	token, err := ioutil.ReadFile("Authentication/BotToken")
	telegramBot.bot, err = tb.NewBot(tb.Settings{
		Token: string(token),
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}
	telegramBot.keyboards = make(map[string]*tb.ReplyMarkup)
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

func (telegramBot *TelegramBot) StartBot() {
	telegramBot.bot.Start()
}

func SendMessage(telegramBotHandler *TelegramBot, usr *User, title string, message interface{}) {

	_, err := telegramBotHandler.bot.Send(usr, title, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}

func (telegramBot *TelegramBot) UserEvent(event interface{}, title string, payload string, responseType uint8) {

	switch responseType {
	case KBOARD:
		telegramBot.bot.Handle(event, func(m *tb.Message) {
			if !m.Private() {
				return
			}
			SendMessage(telegramBot, &me, title, telegramBot.keyboards[payload])
		})
	case TXT:
		telegramBot.bot.Handle(event, func(m *tb.Message) {
			if !m.Private() {
				return
			}
			SendMessage(telegramBot, &me, title, payload)
		})
	}
}
