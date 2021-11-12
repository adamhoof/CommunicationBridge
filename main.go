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

	services.botHandler.CreateBot()

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

	menuKeyboards.AllToys(&telegramBot)
	menuKeyboards.OfficeToys(&telegramBot)

	officeLamp := OfficeLamp{}
	officeCeilLight := OfficeCeilLight{}
	cryptoQuery := CryptoQuery{}
	SetupPhysicalToyInterface(&officeLamp, &services)
	SetupPhysicalToyInterface(&officeCeilLight, &services)
	SetupVirtualToyInterface(&cryptoQuery, &services)
	telegramBot.StartBot()
}
