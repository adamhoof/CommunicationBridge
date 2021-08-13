package CommandDistributor

import (
	"RPICommandHandler/MQTTHandler"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TableLampPub = "room/tableLamp/rpiSet"
)

const (
	TableLampWhiteUpdate = `{"Mode": "white"}`
	TableLampOrangeUpdate = `{"Mode": "orange"}`
	TableLampOffUpdate = `{"Mode": "off"}`
)

func DistributeCommands(mqttClient mqtt.Client, message string) {

	var topic string
	var update interface{}

	switch message {
	case "/tablewhite":
		topic = TableLampPub
		update = TableLampWhiteUpdate
	case "/tableorange":
		topic = TableLampPub
		update = TableLampOrangeUpdate
	case "/tableoff":
		topic = TableLampPub
		update = TableLampOffUpdate
	default:
		return
	}

	MQTTHandler.PublishUpdate(mqttClient, topic, update)
}
