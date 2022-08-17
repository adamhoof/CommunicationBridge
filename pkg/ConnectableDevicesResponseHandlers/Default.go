package deviceresponse

import (
	telegram "CommunicationBridge/pkg/Telegram"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Default(botHandler *telegram.BotHandler, deviceName string) (handler mqtt.MessageHandler) {
	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			botHandler.SendTextMessage(&botHandler.Owner, fmt.Sprintf("%s: %s", deviceName, telegram.ToyCommandIconPairs[msg]))
		}()
	}
	return handler
}
