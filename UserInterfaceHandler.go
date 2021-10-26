package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	APPLIANCES_COMMAND         = "/appliances"
	ALL_APPLIANCES_KEYBOARD    = "allAppliances"
	OFFICE_APPLIANCES_COMMAND  = "/officeappliances"
	OFFICE_APPLIANCES_KEYBOARD = "officeAppliances"
	TABLE_LAMP_COMMAND         = "/tablelamp"
	TABLE_LAMP_KEYBOARD        = "tableLamp"
)

type ApplianceInteractionHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler)
	UserEvents(telegramBot *TelegramBot, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateKeyboard(telegramBot *TelegramBot) map[string]*tb.Btn
	KeyboardRequestHandler(telegramBot *TelegramBot)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, telegramBot *TelegramBot,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()

	keyboard := applianceInteractionHandler.GenerateKeyboard(telegramBot)
	applianceInteractionHandler.UserEvents(telegramBot, mqttHandler, keyboard)
	applianceInteractionHandler.KeyboardRequestHandler(telegramBot)
}