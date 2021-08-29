package main

import (
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}
	telegramBotHandler := TelegramBotHandler{}

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		telegramBotHandler.CreateBot()
		buttons := telegramBotHandler.GenerateButtons()
		telegramBotHandler.TableLampActionHandlers(&mqttHandler, buttons)
		telegramBotHandler.StartBot()
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
}
