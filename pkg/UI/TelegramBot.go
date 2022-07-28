package telegrambot

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"io/ioutil"
	"strings"
	"time"
)

type Handler struct {
	bot              *tb.Bot
	keyboards        map[string]*tb.ReplyMarkup
	lastCommand      string
	eventTimeSetting bool
	eventSetting     bool
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

func (handler *Handler) CreateBot() {

	token, err := ioutil.ReadFile("Auth/BotToken")
	formattedToken := strings.Split(string(token), "\n")
	handler.bot, err = tb.NewBot(tb.Settings{
		Token: formattedToken[0],
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("bot created")
	handler.keyboards = make(map[string]*tb.ReplyMarkup)
}

func CreateHumanReadable(toyDataMap map[string]interface{}) string {

	var humanReadable string

	if toyDataMap != nil {

		for key, value := range toyDataMap {
			switch value.(type) {
			case string:
				humanReadable += key + ": " + value.(string) + "\n"
			case float64:
				value = fmt.Sprintf("%.2f", value.(float64))
				humanReadable += key + ": " + value.(string) + "\n"
			default:
				fmt.Println("unsupported type")
			}
		}
		return humanReadable
	}
	return "empty map"
}

func (handler *Handler) StartBot() {
	handler.bot.Start()
}

func (handler *Handler) SendMessage(telegramBotHandler *Handler, usr *User, title string, message interface{}) {

	_, err := telegramBotHandler.bot.Send(usr, title, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}

func (handler *Handler) UserEvent(event interface{}, title string, payload string, responseType uint8) {

	switch responseType {
	case KBOARD:
		handler.bot.Handle(event, func(c tb.Context) (err error) {
			if !c.Message().Private() {
				return tb.ErrBadRecipient
			}
			handler.SendMessage(handler, &me, title, handler.keyboards[payload])
			return err
		})
	case TXT:
		handler.bot.Handle(event, func(c tb.Context) (err error) {
			if !c.Message().Private() {
				return tb.ErrBadRecipient
			}
			handler.SendMessage(handler, &me, title, payload)
			return err
		})
	}
}
