package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type TableLampActionsHandler struct{}

func (tableLampActionsHandler *TableLampActionsHandler) Name() string{
	return "tableLamp"
}

func (tableLampActionsHandler *TableLampActionsHandler) GenerateButtons(telegramBotHandler *TelegramBotHandler) map[string]*tb.Btn {

	tableLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	tableLampModesMap := make(map[string]*tb.Btn)

	tableLampModesMap["white"] = &tb.Btn{Unique: "white", Text: "⬜"}
	tableLampModesMap["yellow"] = &tb.Btn{Unique: "yellow", Text: "\U0001F7E8"}
	tableLampModesMap["blue"] = &tb.Btn{Unique: "blue", Text: "\U0001F7E6"}
	tableLampModesMap["green"] = &tb.Btn{Unique: "green", Text: "\U0001F7E9"}
	tableLampModesMap["red"] = &tb.Btn{Unique: "red", Text: "\U0001F7E5"}
	tableLampModesMap["pink"] = &tb.Btn{Unique: "pink", Text: "\U0001F7EA"}
	tableLampModesMap["off"] = &tb.Btn{Unique: "off", Text: "🚫"}

	tableLampModesKeyboard.Inline(
		tableLampModesKeyboard.Row(*tableLampModesMap["white"],
			*tableLampModesMap["yellow"], *tableLampModesMap["blue"],
			*tableLampModesMap["green"], *tableLampModesMap["red"],
			*tableLampModesMap["pink"], *tableLampModesMap["off"]),
	)
	telegramBotHandler.keyboards["tableLamp"] = tableLampModesKeyboard
	return tableLampModesMap
}

func (tableLampActionsHandler *TableLampActionsHandler) MessageProcessor() (TableLampMessageHandler mqtt.MessageHandler) {

	TableLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

		tableLampData := make(map[string]interface{})
		tableLampData["Type"] = "TableLamp"
		tableLampData["Mode"] = string(message.Payload())

		func() {
			postgreSQLHandler := PostgreSQLHandler{}
			postgreSQLHandler.Connect()
			postgreSQLHandler.UpdateMode(tableLampData)
			postgreSQLHandler.Disconnect()
		}()
	}
	return TableLampMessageHandler
}

func (tableLampActionsHandler *TableLampActionsHandler) KeyboardRequestHandler(botHandler *TelegramBotHandler) {
	Bot.Handle("/tablelamp", func(message *tb.Message) {
		if !message.Private() {
			return
		}
		_, err := Bot.Send(message.Sender, "Table Lamp Modes", botHandler.keyboards["tableLamp"])
		if err != nil {
			panic(err)
		}
	})
}

func (tableLampActionsHandler *TableLampActionsHandler) SetKeyboardActions(mqttHandler *MQTTHandler, buttons map[string]*tb.Btn) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			Bot.Handle(btn, func(c *tb.Callback) {
				err := Bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				mqttHandler.PublishUpdate(TableLampPub, color)
			})
		}(btn, color)
	}
}
