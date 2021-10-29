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

func (tableLampActionsHandler *TableLampActionsHandler) GenerateKeyboard(telegramBot *TelegramBot) map[string]*tb.Btn {

	tableLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	tableLampModesMap := make(map[string]*tb.Btn)

	tableLampModesMap[white] = &tb.Btn{Unique: white, Text: "⬜"}
	tableLampModesMap[yellow] = &tb.Btn{Unique: yellow, Text: "\U0001F7E8"}
	tableLampModesMap[blue] = &tb.Btn{Unique: blue, Text: "\U0001F7E6"}
	tableLampModesMap[green] = &tb.Btn{Unique: green, Text: "\U0001F7E9"}
	tableLampModesMap[red] = &tb.Btn{Unique: red, Text: "\U0001F7E5"}
	tableLampModesMap[pink] = &tb.Btn{Unique: pink, Text: "\U0001F7EA"}
	tableLampModesMap[off] = &tb.Btn{Unique: off, Text: "🚫"}

	tableLampModesKeyboard.Inline(
		tableLampModesKeyboard.Row(*tableLampModesMap[white],
			*tableLampModesMap[yellow], *tableLampModesMap[blue],
			*tableLampModesMap[green], *tableLampModesMap[red],
			*tableLampModesMap[pink], *tableLampModesMap[off]),
	)
	telegramBot.keyboards[TABLE_LAMP_KEYBOARD] = tableLampModesKeyboard
	return tableLampModesMap
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

func (tableLampActionsHandler *TableLampActionsHandler) KeyboardRequestHandler(telegramBot *TelegramBot) {

	telegramBot.UserEvent(TABLE_LAMP_COMMAND, "Table lamp modes", TABLE_LAMP_KEYBOARD, KBOARD)
}

func (tableLampActionsHandler *TableLampActionsHandler) UserEvents(telegramBot *TelegramBot, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			telegramBot.bot.Handle(btn, func(c *tb.Callback) {
				err := telegramBot.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				mqttHandler.PublishUpdate(tableLampPub, color)
			})
		}(btn, color)
	}
}
