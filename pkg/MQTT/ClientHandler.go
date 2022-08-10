package mqtt2

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type Handler struct {
	client mqtt.Client
}

func (handler *Handler) CreateClient(options *mqtt.ClientOptions) {
	handler.client = mqtt.NewClient(options)
}

func (handler *Handler) Connect() mqtt.Token {
	token := handler.client.Connect()

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}
	fmt.Println("mqtt handler connected")
	return token
}

func (handler *Handler) IsConnected() bool {
	return handler.client.IsConnected()
}

func (handler *Handler) IsConnectionOpen() bool {
	return handler.client.IsConnectionOpen()
}

func (handler *Handler) Disconnect(quiesce uint) {
	handler.client.Disconnect(quiesce)
}

func (handler *Handler) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	token := handler.client.Publish(topic, qos, retained, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (handler *Handler) PublishText(topic string, payload string) mqtt.Token {
	token := handler.client.Publish(topic, 0, true, payload)

	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to send upd: %v", token.Error())
	}
	return token
}

func (handler *Handler) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	token := handler.client.Subscribe(topic, 0, callback)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create subscription: %v", token.Error())
	}
	return token
}

func (handler *Handler) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return handler.client.SubscribeMultiple(filters, callback)
}

func (handler *Handler) Unsubscribe(topics ...string) mqtt.Token {
	return handler.client.Unsubscribe(topics...)
}

func (handler *Handler) AddRoute(topic string, callback mqtt.MessageHandler) {
	handler.client.AddRoute(topic, callback)
}
func (handler *Handler) OptionsReader() mqtt.ClientOptionsReader {
	return handler.client.OptionsReader()
}
