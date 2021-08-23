package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
)

const (
	tableLampSub       = "room/tableLamp/espReply"
)

type MQTTHandler struct {
	tlsConf* tls.Config
	clientOptions mqtt.ClientOptions
	client mqtt.Client
}

func (mqttHandler* MQTTHandler) SetupTLSConfig(){
	certPool := x509.NewCertPool()
	pemCert, certReadingErr := ioutil.ReadFile("Certs/AmazonRootCA1.pem")
	if certReadingErr != nil {
		panic(certReadingErr)
	}
	certPool.AppendCertsFromPEM(pemCert)

	certificateKeyPair, certReadingErr := tls.LoadX509KeyPair("Certs/a29e26a3d1-certificate.pem.crt", "Certs/a29e26a3d1-private.pem.key")
	if certReadingErr != nil {
		panic(certReadingErr)
	}

	mqttHandler.tlsConf = &tls.Config{

		RootCAs: certPool,

		ClientAuth: tls.NoClientCert,

		ClientCAs: nil,

		Certificates: []tls.Certificate{certificateKeyPair},
	}
}

func (mqttHandler* MQTTHandler) SetupClientOptions() {

	mqttHandler.clientOptions.AddBroker("tls://a2z5u1bu7d1g4v-ats.iot.eu-west-2.amazonaws.com:8883")
	mqttHandler.clientOptions.SetClientID("RPICommandHandler").SetTLSConfig(mqttHandler.tlsConf)
	mqttHandler.clientOptions.SetAutoReconnect(true)
	mqttHandler.clientOptions.SetConnectRetry(true)
	mqttHandler.clientOptions.SetCleanSession(true)
	mqttHandler.clientOptions.SetOrderMatters(false)
}

func (mqttHandler *MQTTHandler) CreateClient() {
	mqttHandler.client = mqtt.NewClient(&mqttHandler.clientOptions)
}

func (mqttHandler* MQTTHandler) SetSubscriptions() {
	if token := (mqttHandler.client).Subscribe(tableLampSub, 0, tableLampMessageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	if token := (mqttHandler.client).Subscribe("fuck/shit", 0, tableLampOnBootHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func (mqttHandler* MQTTHandler) ConnectClient() {
	if token := (mqttHandler.client).Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("Client started")
}

var tableLampMessageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {

	applianceDataMap := Collect(message)

	go func() {
		humanReadable := CreateHumanReadable(applianceDataMap)
		userReply := CreateUserReply(humanReadable)
		_, err := Bot.Send(userReply)
		if err != nil {
			panic(err)
		}
	}()


	if applianceDataMap["Mode"] == "failed to set" || applianceDataMap["Mode"] == "already set"{
		return
	}

	go func() {
		db := ConnectDB()
		UpdateMode(db, applianceDataMap)
		CloseConnection(db)
	}()
}

var tableLampOnBootHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	/*db:= ConnectDB()

	message.Payload()
	query := QueryModeProp(db, "TableLamp")
	go CloseConnection(db)

	var update string

	switch query {
	case "white":
		update = TableLampWhiteUpdate
	case "yellow":
		update = TableLampYellowUpdate
	case "red":
		update = TableLampRedUpdate
	case "off":
		update = TableLampOffUpdate
	}
	PublishUpdate(client, tableLampPub, update)*/
	fmt.Println(message.Topic())
	fmt.Println(message.Payload())
}

func PublishUpdate(mqttClient* mqtt.Client, topic string, interfacou interface{}) {

	if token := (*mqttClient).Publish(topic, 0, false, interfacou); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
