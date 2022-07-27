package mqtts

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"log"
)

type Client struct {
	clientOptions mqtt.ClientOptions
	client        mqtt.Client
}

const (
	brokerName = "tls://proteccmqtt.medunka.cz:8883"
	clientName = "RPICommandHandler"
)

func (client *Client) SetupClientOptions() {
	client.clientOptions.AddBroker(brokerName)
	client.clientOptions.SetClientID(clientName)
	client.clientOptions.SetTLSConfig(client.GenerateTlsConfig())
	client.clientOptions.SetAutoReconnect(true)
	client.clientOptions.SetConnectRetry(true)
	client.clientOptions.SetCleanSession(false)
	client.clientOptions.SetOrderMatters(false)
}

func (client *Client) GenerateTlsConfig() *tls.Config {
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

func (client *Client) CreateClient() {
	client.client = mqtt.NewClient(&client.clientOptions)
}

func (client *Client) SetSubscription(messageProcessor mqtt.MessageHandler, topic string) {

	if token := (client.client).Subscribe(topic, 0, messageProcessor); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
}

func (client *Client) ConnectClient() {
	if token := (client.client).Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt client connected")
}

func (client *Client) PublishText(topic string, payload string) {

	if token := (client.client).Publish(topic, 0, true, payload); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
}
