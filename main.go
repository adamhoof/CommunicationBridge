package main

import (
	_ "github.com/lib/pq"
	"sync"
	"time"
)

func main() {
	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}
	telegramBotHandler := TelegramBotHandler{}

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		telegramBotHandler.CreateBot()
		buttons := telegramBotHandler.GenerateButtons()
		telegramBotHandler.TableLampActionHandlers(&mqttHandler, buttons)
		telegramBotHandler.StartBot()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.Disconnect()
	}(&routineSyncer)

	time.Sleep(time.Millisecond*200)
	func() {
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions()
	}()

	routineSyncer.Wait()
}
