package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

type ToyAttributes struct {
	name           string
	availableModes []string
	id             int
	publishTopic   string
	subscribeTopic string
}

var toyColors = map[string]string{
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

func (toyAttributes *ToyAttributes) Name() string {
	return toyAttributes.name
}

func (toyAttributes *ToyAttributes) PubTopic() string {
	return toyAttributes.publishTopic
}

func (toyAttributes *ToyAttributes) SubTopic() string {
	return toyAttributes.subscribeTopic
}

func (toyAttributes *ToyAttributes) MQTTCommandHandler(services *ServiceContainer) {

	handler := func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(toyAttributes.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, toyAttributes.Name()+": "+msg)
			if err != nil {
				return
			}

		}()
	}
	services.mqtt.SetSubscription(handler, toyAttributes.SubTopic())
}

func (toyAttributes *ToyAttributes) GenerateButtons() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	for _, command := range toyAttributes.availableModes {
		func() {
			buttons[command] = &tb.Btn{Unique: command + strconv.Itoa(toyAttributes.id), Text: toyColors[command]}
		}()
	}

	return buttons
}

func (toyAttributes *ToyAttributes) Keyboard(services *ServiceContainer) {

	buttons := toyAttributes.GenerateButtons()
	var buttonsSlice = make([]tb.Btn, len(buttons))

	i := 0
	for name, _ := range buttons {
		buttonsSlice[i] = *buttons[name]
		i++
	}

	keyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	keyboard.Inline(
		keyboard.Row(buttonsSlice...))

	toyAttributes.AwakenButtons(buttons, services)

	services.botHandler.keyboards[toyAttributes.Name()] = keyboard
}

func (toyAttributes *ToyAttributes) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for mode, btn := range buttons {

		func(btn *tb.Btn, mode string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(toyAttributes.PubTopic(), mode)
			})
		}(btn, mode)
	}
}
