package main

type Toy struct {
	name              string
	id                int
	availableCommands []string
	lastKnownCommand  string
	publishTopic      string
	subscribeTopic    string
	keyboardName      string
}
