package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"sync"
)

type MQTTHandler struct {
	tlsConf       *tls.Config
	clientOptions mqtt.ClientOptions
	client        mqtt.Client
	bot           *tb.Bot
}

const (
	TableLampPub = "room/tableLamp/rpiSet"

	tableLampSub = "room/tableLamp/espReply"
)

func (mqttHandler *MQTTHandler) SetupTLSConfig() {
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

func (mqttHandler *MQTTHandler) SetupClientOptions() {

	mqttHandler.clientOptions.AddBroker("tls://a2z5u1bu7d1g4v-ats.iot.eu-west-2.amazonaws.com:8883")
	mqttHandler.clientOptions.SetClientID("RPICommandHandler").SetTLSConfig(mqttHandler.tlsConf)
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
		go func(routineSyncer *sync.WaitGroup) {
			defer routineSyncer.Done()
			humanReadable := CreateHumanReadable(tableLampData)
			SendMessage(mqttHandler.bot, humanReadable, me)
		}(&routineSyncer)

		routineSyncer.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if tableLampData["Mode"] == "already set/failed" {
				return
			}
			postgreSQLHandler := PostgreSQLHandler{}
			postgreSQLHandler.Connect()
			postgreSQLHandler.UpdateMode(tableLampData)
			postgreSQLHandler.CloseConnection()
		}(&routineSyncer)

		routineSyncer.Wait()
	}
	return tableLampMessageHandler
}

func (mqttHandler *MQTTHandler) PublishUpdate(topic string, interfacou interface{}) {

	if token := (mqttHandler.client).Publish(topic, 0, false, interfacou); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
