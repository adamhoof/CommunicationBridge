package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

func GenerateTableLampButtons() (tableLampModes *tb.ReplyMarkup) {
	tableLampModes = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	tableLampOptsMap := make(map[string]*tb.Btn)

	tableLampOptsMap["white"] = &tb.Btn{Unique: "white", Text: "⬜"}
	tableLampOptsMap["yellow"] = &tb.Btn{Unique: "yellow", Text: "\U0001F7E8"}
	tableLampOptsMap["blue"] = &tb.Btn{Unique: "blue", Text: "\U0001F7E6"}
	tableLampOptsMap["green"] = &tb.Btn{Unique: "green", Text: "\U0001F7E9"}
	tableLampOptsMap["red"] = &tb.Btn{Unique: "red", Text: "\U0001F7E5"}
	tableLampOptsMap["pink"] = &tb.Btn{Unique: "pink", Text: "\U0001F7EA"}
	tableLampOptsMap["off"] = &tb.Btn{Unique: "off", Text: "🚫"}

	tableLampModes.Inline(
		tableLampModes.Row(*tableLampOptsMap["white"],
			*tableLampOptsMap["yellow"], *tableLampOptsMap["blue"],
			*tableLampOptsMap["green"], *tableLampOptsMap["red"],
			*tableLampOptsMap["pink"], *tableLampOptsMap["off"]),
	)

	return tableLampModes
}

func TableLampMessageProcessor() (TableLampMessageHandler mqtt.MessageHandler) {

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
