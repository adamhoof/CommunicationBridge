package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TableLampPub = "room/tableLamp/rpiSet"
)

const (
	TableLampWhiteUpdate  = `{"Mode": "white"}`
	TableLampYellowUpdate = `{"Mode": "yellow"}`
	TableLampRedUpdate    = `{"Mode": "red"}`
	TableLampOffUpdate    = `{"Mode": "off"}`
)

func DistributeCommands(mqttClient* mqtt.Client, message string) {

	var topic string
	var update interface{}

	switch message {
	case "/tablewhite":
		topic = TableLampPub
		update = TableLampWhiteUpdate
	case "/tableyellow":
		topic = TableLampPub
		update = TableLampYellowUpdate
	case "/tablered":
		topic = TableLampPub
		update = TableLampRedUpdate
	case "/tableoff":
		topic = TableLampPub
		update = TableLampOffUpdate
	default:
		return
	}

	PublishUpdate(mqttClient, topic, update)
}
