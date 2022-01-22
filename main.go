package main

import (
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler{}
	telegramBot := TelegramBot{}
	postgreSQLHandler := PostgreSQLHandler{}

	services := ServiceContainer{
		mqtt:       &mqttHandler,
		botHandler: &telegramBot,
		db:         &postgreSQLHandler,
	}

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		services.mqtt.SetupClientOptions()
		services.mqtt.CreateClient()
		services.mqtt.ConnectClient()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		services.db.Connect()
		services.db.TestConnection()
	}(&routineSyncer)

	routineSyncer.Wait()

	menuKeyboards := MenuKeyboards{}

	services.botHandler.CreateBot()

	menuKeyboards.AllToys(&telegramBot)
	menuKeyboards.OfficeToys(&telegramBot)
	menuKeyboards.BedroomToys(&telegramBot)

	toyBag := postgreSQLHandler.PullToyData()

	for _, toy := range toyBag {
		toyBag[toy.Name()].MQTTCommandHandler(&services)
		toyBag[toy.Name()].Keyboard(&services)
	}

	telegramBot.StartBot()
}
