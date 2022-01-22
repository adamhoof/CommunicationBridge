package main

type Toy interface {
	Name() string
	PubTopic() string
	SubTopic() string
	MQTTCommandHandler(services *ServiceContainer)
	Keyboard(services *ServiceContainer)
}
