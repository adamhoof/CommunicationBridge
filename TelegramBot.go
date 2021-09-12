package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"sync"
	"time"
)

var Bot* tb.Bot

type TelegramBotHandler struct {

	tableLampModeKeyboard *tb.ReplyMarkup
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
}

func (botHandler *TelegramBotHandler) GenerateButtons() map[string]*tb.Btn {

	tableLampModes := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	m := make(map[string]*tb.Btn)

	m["white"] = &tb.Btn{Unique: "white", Text: "⬜"}
	m["yellow"] = &tb.Btn{Unique: "yellow", Text: "\U0001F7E8"}
	m["blue"] = &tb.Btn{Unique: "blue", Text: "\U0001F7E6"}
	m["green"] = &tb.Btn{Unique: "green", Text: "\U0001F7E9"}
	m["red"] = &tb.Btn{Unique: "red", Text: "\U0001F7E5"}
	m["pink"] = &tb.Btn{Unique: "pink", Text: "\U0001F7EA"}
	m["off"] = &tb.Btn{Unique: "off", Text: "🚫"}

	tableLampModes.Inline(
		tableLampModes.Row(*m["white"], *m["yellow"], *m["blue"], *m["green"], *m["red"], *m["pink"], *m["off"]),
	)
	botHandler.tableLampModeKeyboard = tableLampModes
	return m
}

func (botHandler *TelegramBotHandler) TableLampKeyboardRequestHandler() {
	Bot.Handle("/tablelamp", func(message *tb.Message) {
		if !message.Private() {
			return
		}
		_, err := Bot.Send(message.Sender, "Table Lamp Modes", botHandler.tableLampModeKeyboard)
		if err != nil {
			panic(err)
		}
	})
}

func (botHandler *TelegramBotHandler) TableLampActionHandlers(mqttHandler *MQTTHandler, buttons map[string]*tb.Btn) {
	botHandler.TableLampKeyboardRequestHandler()

	var routineSyncer sync.WaitGroup

	for color, btn := range buttons {

		routineSyncer.Add(1)
		go func(color string, btn *tb.Btn, routineSyncer *sync.WaitGroup) {
			defer routineSyncer.Done()
			Bot.Handle(btn, func(c *tb.Callback) {
				err := Bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				mqttHandler.PublishUpdate(TableLampPub, color)
			})
		}(color, btn, &routineSyncer)
	}
	routineSyncer.Wait()
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

func SendMessage(message string, usr User) {
	_, err := Bot.Send(&usr, message)
	if err != nil {
		panic(err)
	}
}
