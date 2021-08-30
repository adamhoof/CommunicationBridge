package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ProcessJsonData(inputMessage mqtt.Message) (applianceData map[string]interface{}) {

	var rawApplianceData interface{}

	json.Unmarshal(inputMessage.Payload(), &rawApplianceData)

	return rawApplianceData.(map[string]interface{})
}
