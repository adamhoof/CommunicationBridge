package toyresponse

import (
	telegram "RPICommandHandler/pkg/Telegram"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func DefaultDeviceResponseHandler(botHandler *telegram.BotHandler, deviceName string) (handler mqtt.MessageHandler) {
	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			botHandler.SendTextMessage(&botHandler.Owner, fmt.Sprintf("%s: %s", deviceName, msg))
		}()
	}
	return handler
}
