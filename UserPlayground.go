package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PhysicalToy interface {
	Name() string
	MQTTMessageProcessor(services *ServiceContainer) (MessageHandler mqtt.MessageHandler, topic string)
	Kboard(services *ServiceContainer)
	TextCommands(services *ServiceContainer)
}

type VirtualToy interface {
	KeyboardCommands(services *ServiceContainer)
	NonKeyboardCommands(services *ServiceContainer)
}

func SetupPhysicalToyInterface(physicalToy PhysicalToy, services *ServiceContainer) {

	processor, topic := physicalToy.MQTTMessageProcessor(services)
	services.mqtt.SetSubscription(processor, topic)

	physicalToy.Kboard(services)
	physicalToy.TextCommands(services)

	services.db.CreateToy(physicalToy.Name(), "")
}

func SetupVirtualToyInterface(virtualToy VirtualToy, services *ServiceContainer) {
	virtualToy.KeyboardCommands(services)
	virtualToy.NonKeyboardCommands(services)
}
