package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Playground struct {
}

type Toy interface {
	Name() string
	MQTTCommandHandler(services *ServiceContainer) (MessageHandler mqtt.MessageHandler, topic string)
	Kboard(services *ServiceContainer)
	TextCommands(services *ServiceContainer)
}

func (playground *Playground) AddToy(toy Toy, services *ServiceContainer) {

	handler, topic := toy.MQTTCommandHandler(services)
	services.mqtt.SetSubscription(handler, topic)

	toy.Kboard(services)
	toy.TextCommands(services)

	services.db.CreateToy(toy.Name(), "")
}
