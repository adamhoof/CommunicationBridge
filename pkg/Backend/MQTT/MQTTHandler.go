package MQTT

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
)

type Handler struct {
	clientOptions mqtt.ClientOptions
	client        mqtt.Client
}

const (
	brokerName = "tls://proteccmqtt.medunka.cz:8883"
	clientName = "RPICommandHandler"
)

func (handler *Handler) SetupClientOptions() {
	handler.clientOptions.AddBroker(brokerName)
	handler.clientOptions.SetClientID(clientName)
	handler.clientOptions.SetTLSConfig(handler.GenerateTlsConfig())
	handler.clientOptions.SetAutoReconnect(true)
	handler.clientOptions.SetConnectRetry(true)
	handler.clientOptions.SetCleanSession(false)
	handler.clientOptions.SetOrderMatters(false)
}

func (handler *Handler) GenerateTlsConfig() *tls.Config {
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

func (handler *Handler) CreateClient() {
	handler.client = mqtt.NewClient(&handler.clientOptions)
}

func (handler *Handler) SetSubscription(messageProcessor mqtt.MessageHandler, topic string) {

	if token := (handler.client).Subscribe(topic, 0, messageProcessor); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func (handler *Handler) ConnectClient() {
	if token := (handler.client).Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt client connected")
}

func (handler *Handler) PublishText(topic string, payload string) {

	if token := (handler.client).Publish(topic, 0, true, payload); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
