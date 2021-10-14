package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"sync"
	"time"
)

func main() {
	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}
	telegramBotHandler := TelegramBotHandler{}

	handlers := make(map[string]mqtt.MessageHandler)
	handlers["tableLamp"] = TableLampHandler()

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

	routineSyncer.Add(1)
	go func() {
		defer routineSyncer.Done()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions(handlers)
	}()

	routineSyncer.Wait()
}
