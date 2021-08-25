package main

import (
	"fmt"
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
		fmt.Println("running bot")
		done <- true
	}(done)

	go func(chan bool) {
		mqttHandler.SetupTLSConfig()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions()
		fmt.Println("running mqtt")
		done <- true
	}(done)

	go func(chan bool) {
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.CloseConnection()
		fmt.Println("running post")
		done <- true
	}(done)

	_ = <- done
	_ = <- done
	_ = <- done
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
