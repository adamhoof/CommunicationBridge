package mqtt2

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type MQTTClient struct {
}

func (handler *MQTTClient) SetOptions(options mqtt.ClientOptions) {
	handler.SetOptions(options)
}

func (handler *MQTTClient) Connect() mqtt.Token {
	token := handler.Connect()

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt handler connected")
	return token
}

func (handler *MQTTClient) IsConnected() bool {
	return handler.IsConnected()
}

func (handler *MQTTClient) IsConnectionOpen() bool {
	return handler.IsConnectionOpen()
}

func (handler *MQTTClient) Disconnect(quiesce uint) {
	handler.Disconnect(quiesce)
}

func (handler *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	token := handler.Publish(topic, qos, retained, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (handler *MQTTClient) PublishText(topic string, payload string) mqtt.Token {
	token := handler.Publish(topic, 0, true, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (handler *MQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	token := handler.Subscribe(topic, 0, callback)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	return token
}

func (handler *MQTTClient) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return handler.SubscribeMultiple(filters, callback)
}

func (handler *MQTTClient) Unsubscribe(topics ...string) mqtt.Token {
	return handler.Unsubscribe(topics...)
}

func (handler *MQTTClient) AddRoute(topic string, callback mqtt.MessageHandler) {
	handler.AddRoute(topic, callback)
}
func (handler *MQTTClient) OptionsReader() mqtt.ClientOptionsReader {
	return handler.OptionsReader()
}
