package main

import (
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBot := TelegramBot{}
	tableLampActionsHandler := TableLampActionsHandler{}
	keyboardsController := KeyboardsController{}

	telegramBot.CreateBot()

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
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
	SetupApplianceInteractionHandler(&tableLampActionsHandler, &telegramBot, &mqttHandler)
	telegramBot.StartBot()
}
