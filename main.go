package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBot := TelegramBot{}
	tableLampActionsHandler := TableLampActionsHandler{}
	keyboardsController := KeyboardsController{}

	messageProcessors := make(map[string]mqtt.MessageHandler)

	telegramBot.CreateBot()

	SetupClientInterfaceOptions(&tableLampActionsHandler, &telegramBot, &mqttHandler, messageProcessors)

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

	keyboardsController.AllAppliancesKeyboardHandler(&telegramBot)
	keyboardsController.OfficeAppliancesKeyboardHandler(&telegramBot)
	telegramBot.StartBot()
}
