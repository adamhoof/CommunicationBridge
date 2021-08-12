package ApplianceDataCollector

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Collect(inputMessage mqtt.Message) (applianceData []string) {

	var rawApplianceData interface{}

	json.Unmarshal(inputMessage.Payload(), &rawApplianceData)

	applianceDataMap := rawApplianceData.(map[string]interface{})

	if _, ok := applianceDataMap["Params"].(string); !ok{
		applianceData = append(applianceData, applianceDataMap["Type"].(string), applianceDataMap["Mode"].(string))
	} else {

	}
	return applianceData
}
