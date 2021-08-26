package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	var botUpdateConfig tgbotapi.UpdateConfig
	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		SetupBot()
		botUpdateConfig = CreateUpdateConfig()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		mqttHandler.SetupTLSConfig()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.CloseConnection()
	}(&routineSyncer)

	routineSyncer.Wait()

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
