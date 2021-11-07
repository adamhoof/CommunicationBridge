package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
)

type MQTTHandler struct {
	clientOptions mqtt.ClientOptions
	client        mqtt.Client
}

const (
	brokerName = "tls://proteccmqtt.medunka.cz:8883"
	clientName = "RPICommandHandler"
)

func (mqttHandler *MQTTHandler) SetupClientOptions() {
	mqttHandler.clientOptions.AddBroker(brokerName)
	mqttHandler.clientOptions.SetClientID(clientName)
	mqttHandler.clientOptions.SetTLSConfig(mqttHandler.GenerateTlsConfig())
	mqttHandler.clientOptions.SetAutoReconnect(true)
	mqttHandler.clientOptions.SetConnectRetry(true)
	mqttHandler.clientOptions.SetCleanSession(false)
	mqttHandler.clientOptions.SetOrderMatters(false)
}

func (mqttHandler *MQTTHandler) GenerateTlsConfig() *tls.Config {
	certpool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("Auth/ca/ca.crt")
	if err != nil {
		log.Fatalln(err.Error())
	}
	certpool.AppendCertsFromPEM(ca)

	certificateKeyPair, certReadingErr := tls.LoadX509KeyPair("Auth/client/client.crt", "Auth/client/client.key")

	if certReadingErr != nil {
		panic(certReadingErr)
	}

	return &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		Certificates:       []tls.Certificate{certificateKeyPair},
	}
}

func (mqttHandler *MQTTHandler) CreateClient() {
	mqttHandler.client = mqtt.NewClient(&mqttHandler.clientOptions)
}

func (mqttHandler *MQTTHandler) SetSubscription(messageProcessor mqtt.MessageHandler, topic string) {

	if token := (mqttHandler.client).Subscribe(topic, 0, messageProcessor); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func (mqttHandler *MQTTHandler) ConnectClient() {
	if token := (mqttHandler.client).Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt client connected")
}

func (mqttHandler *MQTTHandler) PublishText(topic string, payload string) {

	if token := (mqttHandler.client).Publish(topic, 0, true, payload); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
