package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"sync"
)

type MQTTHandler struct {
	clientOptions mqtt.ClientOptions
	client        mqtt.Client
}

const (
	TableLampPub = "room/tableLamp/rpiSet"

	tableLampSub = "room/tableLamp/espReply"
)

func (mqttHandler *MQTTHandler) SetupClientOptions() {

	mqttHandler.clientOptions.AddBroker("tcp://185.152.65.53:1883")
	mqttHandler.clientOptions.SetClientID("RPICommandHandler")
	mqttHandler.clientOptions.SetUsername("device")
	mqttHandler.clientOptions.SetPassword("devicepasswrd")
	mqttHandler.clientOptions.SetAutoReconnect(true)
	mqttHandler.clientOptions.SetConnectRetry(true)
	mqttHandler.clientOptions.SetCleanSession(false)
	mqttHandler.clientOptions.SetOrderMatters(false)
}

func (mqttHandler *MQTTHandler) CreateClient() {
	mqttHandler.client = mqtt.NewClient(&mqttHandler.clientOptions)
}

func (mqttHandler *MQTTHandler) SetSubscriptions() {
	if token := (mqttHandler.client).Subscribe(tableLampSub, 0, mqttHandler.TableLampHandler()); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func (mqttHandler *MQTTHandler) ConnectClient() {
	if token := (mqttHandler.client).Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("Client started")
}

func (mqttHandler *MQTTHandler) TableLampHandler() (tableLampMessageHandler mqtt.MessageHandler) {

	tableLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

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
	return tableLampMessageHandler
}

func (mqttHandler *MQTTHandler) PublishUpdate(topic string, interfacou interface{}) {

	if token := (mqttHandler.client).Publish(topic, 0, true, interfacou); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
