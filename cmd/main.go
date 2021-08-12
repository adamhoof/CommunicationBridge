package main

import (
	"RPICommandHandler/CommandDistributor"
	"RPICommandHandler/MQTTHandler"
	"RPICommandHandler/PostgreSQLHandler"
	"RPICommandHandler/TelegramBot"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {

	TelegramBot.Setup()

	db := PostgreSQLHandler.Connect()

	PostgreSQLHandler.TestConnection(db)

	tlsConfig := MQTTHandler.NewTLSConfig()

	clientOptions := MQTTHandler.SetupClientOptions(tlsConfig)
	mqttClient := mqtt.NewClient(clientOptions)

	MQTTHandler.ConnectClient(mqttClient)

	MQTTHandler.SetSubscriptions(mqttClient)

	telegramBotUpdate := tgbotapi.NewUpdate(0)
	telegramBotUpdate.Timeout = 60

	updates, err := TelegramBot.Bot.GetUpdatesChan(telegramBotUpdate)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		message := update.Message.Text

		CommandDistributor.DistributeCommands(mqttClient, message)
	}
}
