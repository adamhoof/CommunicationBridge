package main

type Toy struct {
	name              string
	availableCommands []string
	lastKnownCommand  string
	id                int
	publishTopic      string
	subscribeTopic    string
}
