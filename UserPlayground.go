package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Toy interface {
	Name() string
	MQTTProcessor(services *ServiceContainer) (MessageHandler mqtt.MessageHandler, topic string)
	Kboard(services *ServiceContainer)
	TextCommands(services *ServiceContainer)
}

func SetupPhysicalToyInterface(physicalToy Toy, services *ServiceContainer) {

	processor, topic := physicalToy.MQTTProcessor(services)
	services.mqtt.SetSubscription(processor, topic)

	physicalToy.Kboard(services)
	physicalToy.TextCommands(services)

	services.db.CreateToy(physicalToy.Name(), "")
}
