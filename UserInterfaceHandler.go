package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type ApplianceInteractionHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler)
	SetKeyboardActions(telegramBotHandler *TelegramBotHandler, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateKeyboard(telegramBotHandler *TelegramBotHandler) map[string]*tb.Btn
	KeyboardRequestHandler(botHandler *TelegramBotHandler)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, telegramBotHandler *TelegramBotHandler,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()

	keyboard := applianceInteractionHandler.GenerateKeyboard(telegramBotHandler)
	applianceInteractionHandler.SetKeyboardActions(telegramBotHandler, mqttHandler, keyboard)
	applianceInteractionHandler.KeyboardRequestHandler(telegramBotHandler)
}

func RoomAppliancesKeyboardRequestHandler(telegramBotHandler *TelegramBotHandler) {
	roomAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBotHandler.keyboards["roomAppliances"] = roomAppliancesKeyboard

	tableLampBtn := roomAppliancesKeyboard.Text("/tablelamp")
	roomAppliancesKeyboard.Reply(
		roomAppliancesKeyboard.Row(tableLampBtn),
	)
	usr := User{userId: "558297691"}

	telegramBotHandler.bot.Handle("/roomappliances", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		SendMessage(telegramBotHandler, usr, "Room Appliances", telegramBotHandler.keyboards["roomAppliances"])
	})
	telegramBotHandler.bot.Handle(&tableLampBtn, func(m *tb.Message) {
		SendMessage(telegramBotHandler, usr, "Table lamp modes", telegramBotHandler.keyboards["tableLamp"])
	})
}
