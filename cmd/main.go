package main

import (
	"RPICommandHandler/pkg/Backend/Database"
	"RPICommandHandler/pkg/Backend/MQTT"
	telegrambot "RPICommandHandler/pkg/Frontend"
	_ "github.com/lib/pq"
	"sync"
)

func main() {

	mqttHandler := mqtts.Client{}
	telegramBot := telegrambot.Handler{}
	postsgresSQLHandler := database.PostgresSQLHandler{}

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
		postsgresSQLHandler.Connect()
		postsgresSQLHandler.TestConnection()
	}(&routineSyncer)

	routineSyncer.Wait()

	menuKeyboards := MenuKeyboards{}

	telegramBot.CreateBot()

	menuKeyboards.AllToys(&telegramBot)
	menuKeyboards.OfficeToys(&telegramBot)
	menuKeyboards.BedroomToys(&telegramBot)

	toyBag := postsgresSQLHandler.PullToyData()

	for _, toy := range toyBag {
		toyBag[toy.Name()].MQTTCommandHandler(&
		toyBag[toy.Name()].Keyboard(&
	}

	telegramBot.StartBot()
}
