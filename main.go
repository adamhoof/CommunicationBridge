package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBotHandler := TelegramBotHandler{}
	tableLampActionsHandler := TableLampActionsHandler{}

	messageProcessors := make(map[string]mqtt.MessageHandler)

	telegramBotHandler.CreateBot()

	SetupClientInterfaceOptions(&tableLampActionsHandler, &telegramBotHandler, &mqttHandler, messageProcessors)

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions(messageProcessors)
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		postgreSQLHandler := PostgreSQLHandler{}
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.Disconnect()
	}(&routineSyncer)

	routineSyncer.Wait()

	AllAppliancesKeyboardHandler(&telegramBotHandler)
	OfficeAppliancesKeyboardHandler(&telegramBotHandler)
	telegramBotHandler.StartBot()
}
