package main

import (
	_ "github.com/lib/pq"
)

func main() {

	SetupBot()

	mqttHandler := MQTTHandler {}
	mqttHandler.SetupTLSConfig()
	mqttHandler.SetupClientOptions()
	mqttHandler.CreateClient()
	mqttHandler.ConnectClient()
	mqttHandler.SetSubscriptions()

	db := ConnectDB()

	TestDBConnection(db)

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

		DistributeCommands(&mqttHandler.client, message)
	}
}
