package device_action

import (
	connectable "CommunicationBridge/pkg/ConnectableDevices"
	database "CommunicationBridge/pkg/Database"
	telegram "CommunicationBridge/pkg/Telegram"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/telebot.v3"
)

func OnToyBoot(dbHandler database.DatabaseHandler, botHandler *telegram.BotHandler, mqttClient mqtt.Client, keyboards *map[string]*tb.ReplyMarkup) (handler mqtt.MessageHandler) {
	handler = func(client mqtt.Client, message mqtt.Message) {
		var toy connectable.Toy
		err := json.Unmarshal(message.Payload(), &toy)
		if err != nil {
			fmt.Println(err)
		}
		toy.BotCommand = "/" + toy.Name

		dbHandler.RegisterToy(&toy)
		buttons := telegram.GenerateToyButtonsWithClickHandlers(botHandler, mqttClient, &toy)
		keyboard := telegram.GenerateToyKeyboard(buttons)
		(*keyboards)[toy.Name] = keyboard
		botHandler.HandleCommand(toy.BotCommand, botHandler.SendKeyboard(toy.Name, *keyboards, toy.Name))
		mqttClient.Subscribe(toy.SubscribeTopic, 0, Default(botHandler, toy.Name))
	}
	return handler
}
