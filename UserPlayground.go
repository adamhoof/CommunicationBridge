package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Playground struct {
}

type Toy interface {
	Name() string
	PubTopic() string
	SubTopic() string
	MQTTCommandHandler(services *ServiceContainer) (MessageHandler mqtt.MessageHandler)
	Keyboard(services *ServiceContainer)
}
