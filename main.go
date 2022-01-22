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
	menuKeyboards.BedroomToys(&telegramBot)

	playground := Playground{}
	playground.ColorTheToys()
	toyBag := ToyBag{}

	toyBag.bag = make(map[string]Toy)

	postgreSQLHandler.PullToyData(toyBag.bag)

	for _, toy := range toyBag.bag {
		handler, topic := toyBag.bag[toy.Name()].MQTTCommandHandler(&services)
		services.mqtt.SetSubscription(handler, topic)
		toyBag.bag[toy.Name()].Keyboard(&services)
	}

	telegramBot.StartBot()
}
