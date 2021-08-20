package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	botUpdateConfig := CreateUpdateConfig()

	updates, err := Bot.GetUpdatesChan(botUpdateConfig)
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
