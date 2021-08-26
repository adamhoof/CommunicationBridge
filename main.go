package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {

	var botUpdateConfig tgbotapi.UpdateConfig
	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}

	done := make(chan bool)

	go func(chan bool) {
		SetupBot()
		botUpdateConfig = CreateUpdateConfig()
		done <- true
	}(done)

	go func(chan bool) {
		mqttHandler.SetupTLSConfig()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions()
		done <- true
	}(done)

	go func(chan bool) {
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.CloseConnection()
		done <- true
	}(done)

	<- done
	<- done
	<- done
	close(done)

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
