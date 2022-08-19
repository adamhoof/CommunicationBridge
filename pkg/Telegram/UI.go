package telegram

import (
	connectable "CommunicationBridge/pkg/ConnectableDevices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

var CommandPairIcons = map[string]string{
	"on":     "⬜",
	"white":  "⬜",
	"yellow": "\U0001F7E8",
	"blue":   "\U0001F7E6",
	"green":  "\U0001F7E9",
	"red":    "\U0001F7E5",
	"pink":   "\U0001F7EA",
	"orange": "\U0001F7E7",
	"off":    "🚫",
	"open":   "🌞",
	"close":  "🌚"}

func GenerateToyKeyboard(buttons []tb.Btn) *tb.ReplyMarkup {

	keyboard := &tb.ReplyMarkup{
		ResizeKeyboard: true,
	}
	keyboard.Inline(keyboard.Row(buttons...))

	return keyboard
}

func GenerateToyButtonsWithClickHandlers(botHandler *BotHandler, client mqtt.Client, toy *connectable.Toy) (buttons []tb.Btn) {
	func() {
		for _, command := range toy.AvailableModes {
			button := tb.Btn{
				Unique: command + strconv.Itoa(toy.Id),
				Text:   CommandPairIcons[command]}
			buttons = append(buttons, button)

			botHandler.HandleButtonClick(&button, botHandler.SendCommandViaMQTT(command, toy.PublishTopic, client))
		}
	}()
	return buttons
}
