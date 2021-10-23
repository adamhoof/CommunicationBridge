package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type ApplianceInteractionHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler)
	SetKeyboardActions(mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateButtons(telegramBotHandler *TelegramBotHandler) map[string]*tb.Btn
	KeyboardRequestHandler(botHandler *TelegramBotHandler)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, telegramBotHandler *TelegramBotHandler,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()
	applianceInteractionHandler.SetKeyboardActions(mqttHandler,
		applianceInteractionHandler.GenerateButtons(telegramBotHandler))
	applianceInteractionHandler.KeyboardRequestHandler(telegramBotHandler)
}

func AllAppliancesKeyboard(telegramBotHandler *TelegramBotHandler)  {
	allAppliancesKeyboard := &tb.ReplyMarkup{}
	tableLampBtn := allAppliancesKeyboard.Text("/Table \U0001FA94")
	telegramBotHandler.keyboards["appliances"] = allAppliancesKeyboard

	allAppliancesKeyboard.Reply(
		allAppliancesKeyboard.Row(tableLampBtn),
		)
	usr := User{userId: "558297691"}

	Bot.Handle("/appliances", func(m *tb.Message){
		if !m.Private(){
			return
		}
		SendMessage(usr, "Appliances", telegramBotHandler.keyboards["appliances"])
	})
	Bot.Handle(&tableLampBtn, func(m *tb.Message){
		SendMessage(usr, "Table \U0001FA94 modes", telegramBotHandler.keyboards["tableLamp"])
	})
}
