package deviceresponse

import (
	connectable "CommunicationBridge/pkg/ConnectableDevices"
	database "CommunicationBridge/pkg/Database"
	telegram "CommunicationBridge/pkg/Telegram"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/telebot.v3"
)

func OnNewDeviceBoot(dbHandler database.DatabaseHandler, botHandler *telegram.BotHandler, mqttClient mqtt.Client, keyboards *map[string]*tb.ReplyMarkup) (handler mqtt.MessageHandler) {
	handler = func(client mqtt.Client, message mqtt.Message) {
		var toy connectable.Toy
		err := json.Unmarshal(message.Payload(), &toy)
		if err != nil {
			fmt.Println(err)
		}
		toy.BotCommand = "/" + toy.Name

		dbHandler.RegisterToy(&toy)
		keyboard := telegram.GenerateKeyboardWithButtonsHandlersForToy(botHandler, mqttClient, &toy)
		(*keyboards)[toy.Name] = keyboard
		botHandler.HandleCommand(toy.BotCommand, botHandler.SendKeyboard(toy.Name, *keyboards, toy.Name))
		mqttClient.Subscribe(toy.SubscribeTopic, 0, Default(botHandler, toy.Name))
	}
	return handler
}
