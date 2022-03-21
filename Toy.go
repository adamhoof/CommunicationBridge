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

func (toy *Toy) assignKeyboardName(name string) {
	toy.keyboardName = name
}
