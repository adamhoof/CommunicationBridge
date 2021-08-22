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
	tableLampOnBootSub = "room/tableLamp/espOnBoot"
	tableLampSub       = "room/tableLamp/espReply"
)

const (
	tableLampPub = "room/tableLamp/rpiSet"
)

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
	db:= ConnectDB()

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
	PublishUpdate(client, tableLampPub, update)
}

func NewTLSConfig() (config *tls.Config) {

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

	config = &tls.Config{

		RootCAs: certPool,

		ClientAuth: tls.NoClientCert,

		ClientCAs: nil,

		Certificates: []tls.Certificate{certificateKeyPair},
	}
	return
}

func SetupClientOptions(config *tls.Config) (clientOptions *mqtt.ClientOptions) {

	clientOptions = mqtt.NewClientOptions()
	clientOptions.AddBroker("tls://a2z5u1bu7d1g4v-ats.iot.eu-west-2.amazonaws.com:8883")
	clientOptions.SetClientID("RPICommandHandler").SetTLSConfig(config)
	clientOptions.SetAutoReconnect(true)
	clientOptions.SetConnectRetry(true)
	clientOptions.SetCleanSession(false)
	clientOptions.SetOrderMatters(false)

	return clientOptions
}

func ConnectClient(mqttClient mqtt.Client) {
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("Client started")
}

func SetMQTTSubscriptions(mqttClient mqtt.Client) {
	if token := mqttClient.Subscribe(tableLampOnBootSub, 1, tableLampOnBootHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	if token := mqttClient.Subscribe(tableLampSub, 1, tableLampMessageHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func PublishUpdate(mqttClient mqtt.Client, topic string, interfacou interface{}) {

	if token := mqttClient.Publish(topic, 1, false, interfacou); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
