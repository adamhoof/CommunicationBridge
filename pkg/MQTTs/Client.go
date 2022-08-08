package mqtts

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type MQTTClient struct {
}

func (client *MQTTClient) SetOptions(options mqtt.ClientOptions) {
	client.SetOptions(options)
}

func (client *MQTTClient) Connect() mqtt.Token {
	token := client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt handler connected")
	return token
}

func (client *MQTTClient) IsConnected() bool {
	return client.IsConnected()
}

func (client *MQTTClient) IsConnectionOpen() bool {
	return client.IsConnectionOpen()
}

func (client *MQTTClient) Disconnect(quiesce uint) {
	client.Disconnect(quiesce)
}

func (client *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	token := client.Publish(topic, qos, retained, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (client *MQTTClient) PublishText(topic string, payload string) mqtt.Token {
	token := client.Publish(topic, 0, true, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (client *MQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	token := client.Subscribe(topic, 0, callback)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	return token
}

func (client *MQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return client.SubscribeMultiple(filters, callback)
}

func (client *MQTTClient) Unsubscribe(topics ...string) mqtt.Token {
	return client.Unsubscribe(topics...)
}

func (client *MQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {
	client.AddRoute(topic, callback)
}
func (client *MQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	return client.OptionsReader()
}
