package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type ApplianceActionsHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler)
	SetKeyboardActions(mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateButtons(telegramBotHandler *TelegramBotHandler) map[string]*tb.Btn
	KeyboardRequestHandler(botHandler *TelegramBotHandler)
}

func SetupClientInterfaceOptions(clientCommandHandler ApplianceActionsHandler, telegramBotHandler *TelegramBotHandler,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[clientCommandHandler.Name()] = clientCommandHandler.MessageProcessor()
	clientCommandHandler.SetKeyboardActions(mqttHandler,
		clientCommandHandler.GenerateButtons(telegramBotHandler))
	clientCommandHandler.KeyboardRequestHandler(telegramBotHandler)
}
