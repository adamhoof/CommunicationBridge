package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

func TableLampHandler() (TableLampMessageHandler mqtt.MessageHandler) {

	TableLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

		tableLampData := make(map[string]interface{})
		tableLampData["Type"] = "TableLamp"
		tableLampData["Mode"] = string(message.Payload())

		var routineSyncer sync.WaitGroup

		routineSyncer.Add(1)
		go func() {
			defer routineSyncer.Done()
			me := User{userId: "558297691"}
			SendMessage(CreateHumanReadable(tableLampData), me)
		}()

		routineSyncer.Add(1)
		go func() {
			defer routineSyncer.Done()
			postgreSQLHandler := PostgreSQLHandler{}
			postgreSQLHandler.Connect()
			postgreSQLHandler.UpdateMode(tableLampData)
			postgreSQLHandler.Disconnect()
		}()

		routineSyncer.Wait()
	}
	return TableLampMessageHandler
}
