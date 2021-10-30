package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	ALL_APPLIANCES_KEYBOARD    = "allAppliances"
	OFFICE_APPLIANCES_KEYBOARD = "officeAppliances"
	TABLE_LAMP_KEYBOARD        = "tableLamp"
)

type ApplianceInteractionHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler, topic string)
	GenerateKeyboard(telegramBot *TelegramBot, mqttHandler* MQTTHandler)
}

func SetupApplianceInteractionHandler(applianceInteractionHandler ApplianceInteractionHandler, telegramBot *TelegramBot,
	mqttHandler *MQTTHandler) {

	processor, topic := applianceInteractionHandler.MessageProcessor()
	mqttHandler.SetSubscription(processor, topic)

	applianceInteractionHandler.GenerateKeyboard(telegramBot, mqttHandler)
}
