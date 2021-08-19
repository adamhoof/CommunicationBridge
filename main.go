package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {

	SetupBot()

	db := ConnectDB()

	TestDBConnection(db)

	tlsConfig := NewTLSConfig()

	clientOptions := SetupClientOptions(tlsConfig)
	mqttClient := mqtt.NewClient(clientOptions)

	ConnectClient(mqttClient)

	SetMQTTSubscriptions(mqttClient)

	telegramBotUpdate := tgbotapi.NewUpdate(0)
	telegramBotUpdate.Timeout = 60

	updates, err := Bot.GetUpdatesChan(telegramBotUpdate)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

  		message := update.Message.Text

		DistributeCommands(mqttClient, message)
	}
}
