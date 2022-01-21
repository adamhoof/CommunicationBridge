package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Playground struct {
}

type ToyAttributes struct {
	name            string
	commandWithName []byte
	uniqueConst     string
	publishTopic    string
	subscribeTopic  string
}

func (toyAttributes *ToyAttributes) Name() string {
	return toyAttributes.name
}

func (toyAttributes *ToyAttributes) PubTopic() string {
	return toyAttributes.publishTopic
}

func (toyAttributes *ToyAttributes) SubTopic() string {
	return toyAttributes.subscribeTopic
}

func (toyAttributes *ToyAttributes) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(toyAttributes.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, toyAttributes.Name()+": "+msg)
			if err != nil {
				return
			}

		}()
	}
	return handler, toyAttributes.SubTopic()
}

func (toyAttributes *ToyAttributes) GenerateKboardBtns() map[string]*tb.Btn {

	commandWithName := make(map[string]interface{})
	err := json.Unmarshal(toyAttributes.commandWithName, &commandWithName)
	if err != nil {
		fmt.Println("shitter")
	}

	buttons := make(map[string]*tb.Btn)

	for command, name := range commandWithName {
		buttons[command] = &tb.Btn{Unique: command + toyAttributes.uniqueConst, Text: name.(string)}
	}
	/*
		buttons[officeLampWhite] = &tb.Btn{Unique: officeLampWhite + toyAttributes.uniqueConst, Text: "⬜"}
		buttons[officeLampYellow] = &tb.Btn{Unique: officeLampYellow + toyAttributes.uniqueConst, Text: "\U0001F7E8"}
		buttons[officeLampBlue] = &tb.Btn{Unique: officeLampBlue + toyAttributes.uniqueConst, Text: "\U0001F7E6"}
		buttons[officeLampGreen] = &tb.Btn{Unique: officeLampGreen + toyAttributes.uniqueConst, Text: "\U0001F7E9"}
		buttons[officeLampRed] = &tb.Btn{Unique: officeLampRed + toyAttributes.uniqueConst, Text: "\U0001F7E5"}
		buttons[officeLampPink] = &tb.Btn{Unique: officeLampPink + toyAttributes.uniqueConst, Text: "\U0001F7EA"}
		buttons[officeLampOff] = &tb.Btn{Unique: officeLampOff + toyAttributes.uniqueConst, Text: "🚫"}*/

	return buttons
}

func (toyAttributes *ToyAttributes) Kboard(services *ServiceContainer) {

	buttons := toyAttributes.GenerateKboardBtns()
	var buttonsSlice = make([]tb.Btn, len(buttons))

	i := 1
	for name, _ := range buttons {
		buttonsSlice[i] = *buttons[name]
	}

	officeLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(buttonsSlice...))

	/*officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(*buttons[officeLampWhite],
			*buttons[officeLampYellow], *buttons[officeLampBlue],
			*buttons[officeLampGreen], *buttons[officeLampRed],
			*buttons[officeLampPink], *buttons[officeLampOff]),
	)*/

	toyAttributes.AwakenButtons(buttons, services)

	services.botHandler.keyboards[OfficeLampKeyboard] = officeLampModesKeyboard
}

func (toyAttributes *ToyAttributes) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(toyAttributes.PubTopic(), color)
			})
		}(btn, color)
	}
}

func (toyAttributes *ToyAttributes) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officelamp", "Office lamp", OfficeLampKeyboard, KBOARD)
}

type Toy interface {
	Name() string
	PubTopic() string
	SubTopic() string
	MQTTCommandHandler(services *ServiceContainer) (MessageHandler mqtt.MessageHandler, topic string)
	Kboard(services *ServiceContainer)
	TextCommands(services *ServiceContainer)
}

func (playground *Playground) spawnToys(services *ServiceContainer) {

}

func (playground *Playground) takeOutToys(toyStorage *ToyBag, services *ServiceContainer) {

	for _, toy := range toyStorage.bag {
		handler, topic := toy.MQTTCommandHandler(services)
		services.mqtt.SetSubscription(handler, topic)

		toy.Kboard(services)
		toy.TextCommands(services)

		services.db.CreateToy(toy.Name(), "")
	}
}
