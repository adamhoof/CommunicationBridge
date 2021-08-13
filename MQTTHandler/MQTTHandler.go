package MQTTHandler

import (
	"RPICommandHandler/ApplianceDataCollector"
	"RPICommandHandler/PostgreSQLHandler"
	"RPICommandHandler/TelegramBot"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
	"strings"
)

const (
	tableLampOnBootSub = "room/tableLamp/espOnBoot"
	tableLampSub       = "room/tableLamp/espReply"
)

var messageHandler mqtt.MessageHandler = func(mqttClient mqtt.Client, msg mqtt.Message) {

	switch strings.Contains(msg.Topic(), "espOnBoot") {
	case false:

		applianceDataMap := ApplianceDataCollector.Collect(msg)

		humanReadable := TelegramBot.CreateHumanReadable(applianceDataMap)
		userReply := TelegramBot.CreateUserReply(humanReadable)
		go TelegramBot.Bot.Send(userReply)

		if applianceDataMap["Mode"] == "failed to set" || applianceDataMap["Mode"] == "already set"{
			return}

		go func() {
			db := PostgreSQLHandler.Connect()
			PostgreSQLHandler.UpdateMode(db, applianceDataMap)
			PostgreSQLHandler.CloseConnection(db)}()
		return

	case true:
		if strings.Contains(msg.Topic(), "tableLamp") {
			db:= PostgreSQLHandler.Connect()

			query := PostgreSQLHandler.QueryModeProp(db, "TableLamp")
			go PostgreSQLHandler.CloseConnection(db)

			var update string

			switch query {
			case "white":
				update = `{"Mode": "white"}`
			case "orange":
				update = `{"Mode": "orange"}`
			case "off":
				update = `{"Mode": "off"}`
			}
			PublishUpdate(mqttClient, "room/tableLamp/rpiSet", update)
		}
	}
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
	clientOptions.SetDefaultPublishHandler(messageHandler)
	clientOptions.SetAutoReconnect(true)
	clientOptions.SetConnectRetry(true)
	clientOptions.SetCleanSession(false)

	return clientOptions
}

func ConnectClient(mqttClient mqtt.Client) {
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("Client started")
}

func SetSubscriptions(mqttClient mqtt.Client) {

	if token := mqttClient.Subscribe(tableLampOnBootSub, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	if token := mqttClient.Subscribe(tableLampSub, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func PublishUpdate(mqttClient mqtt.Client, topic string, interfacou interface{}) {

	if token := mqttClient.Publish(topic, 1, false, interfacou); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
