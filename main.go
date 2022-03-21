package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"strconv"
	"sync"
)

func defaultResponseHandler(toy *Toy, services *ServiceContainer) func(client mqtt.Client, message mqtt.Message) {

	handler := func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateToyMode(toy.name, toy.lastKnownCommand)
			_, err := services.botHandler.bot.Send(&me, toy.name+": "+toy.lastKnownCommand)
			if err != nil {
				return
			}

		}()
	}
	return handler
}

func main() {

	mqttHandler := MQTTHandler{}
	telegramBot := TelegramBot{}
	dbHandler := DBHandler{}

	services := ServiceContainer{
		mqtt:       &mqttHandler,
		botHandler: &telegramBot,
		db:         &dbHandler,
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

	buttonFactory := ButtonFactory{}
	keyboardFactory := KeyboardFactory{}
	keyboardStorage := KeyboardStorage{}
	keyboardStorage.unlock()

	commandIconTemplates := map[string]string{
		"on":     "⬜",
		"white":  "⬜",
		"yellow": "\U0001F7E8",
		"blue":   "\U0001F7E6",
		"green":  "\U0001F7E9",
		"red":    "\U0001F7E5",
		"pink":   "\U0001F7EA",
		"orange": "\U0001F7E7",
		"off":    "🚫",
		"1":      "🌞",
		"0":      "🌚"}

	buttonFactory.setCommandAndIconButtonTemplates(commandIconTemplates)

	services.botHandler.CreateBot()

	toyBag := dbHandler.PullToyData()

	for _, toy := range toyBag {

		buttons, err := buttonFactory.generateInlineButtons(toy.id, toy.availableCommands)
		if err != nil {
			fmt.Println(err)
		}
		buttonFactory.setButtonHandlers(buttons, toy, &services)

		toy.assignKeyboardName(toy.name + strconv.Itoa(toy.id))
		keyboard := keyboardFactory.createFromButtons(buttons)
		keyboardStorage.store(toy.keyboardName, keyboard)

		mqttHandler.SetSubscription(defaultResponseHandler(toy, &services), toy.subscribeTopic)
	}

	telegramBot.StartBot()
}
