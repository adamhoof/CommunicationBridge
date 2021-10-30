package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	white        = "w"
	yellow       = "y"
	blue         = "b"
	green        = "g"
	red          = "r"
	pink         = "p"
	off          = "o"
	tableLampPub = "room/tableLamp/rpiSet"
	tableLampSub = "room/tableLamp/espReply"
)

type TableLampActionsHandler struct{}

func (tableLampActionsHandler *TableLampActionsHandler) Name() string {
	return "tableLamp"
}

func (tableLampActionsHandler *TableLampActionsHandler) GenerateFunctionButtons(botHandler *TelegramBot, mqttHandler *MQTTHandler) map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[white] = &tb.Btn{Unique: white, Text: "⬜"}
	buttons[yellow] = &tb.Btn{Unique: yellow, Text: "\U0001F7E8"}
	buttons[blue] = &tb.Btn{Unique: blue, Text: "\U0001F7E6"}
	buttons[green] = &tb.Btn{Unique: green, Text: "\U0001F7E9"}
	buttons[red] = &tb.Btn{Unique: red, Text: "\U0001F7E5"}
	buttons[pink] = &tb.Btn{Unique: pink, Text: "\U0001F7EA"}
	buttons[off] = &tb.Btn{Unique: off, Text: "🚫"}

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				mqttHandler.PublishUpdate(tableLampPub, color)
			})
		}(btn, color)
	}
	return buttons
}

func (tableLampActionsHandler *TableLampActionsHandler) GenerateKeyboard(telegramBot *TelegramBot, mqtt *MQTTHandler) {

	buttons := tableLampActionsHandler.GenerateFunctionButtons(telegramBot, mqtt)

	tableLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	tableLampModesKeyboard.Inline(
		tableLampModesKeyboard.Row(*buttons[white],
			*buttons[yellow], *buttons[blue],
			*buttons[green], *buttons[red],
			*buttons[pink], *buttons[off]),
	)
	telegramBot.keyboards[TABLE_LAMP_KEYBOARD] = tableLampModesKeyboard
}

func (tableLampActionsHandler *TableLampActionsHandler) MessageProcessor() (TableLampMessageHandler mqtt.MessageHandler, topic string) {

	TableLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			postgreSQLHandler := PostgreSQLHandler{}
			postgreSQLHandler.Connect()
			postgreSQLHandler.UpdateMode("TableLamp", string(message.Payload()))
			postgreSQLHandler.Disconnect()
		}()
	}
	return TableLampMessageHandler, tableLampSub
}

