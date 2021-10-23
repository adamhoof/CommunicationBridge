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
