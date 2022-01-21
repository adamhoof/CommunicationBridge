package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	officeLampWhite    = "white"
	officeLampYellow   = "yellow"
	officeLampBlue     = "blue"
	officeLampGreen    = "green"
	officeLampRed      = "red"
	officeLampPink     = "pink"
	officeLampOff      = "off"
	OfficeLampKeyboard = "officelamp"
)

type OfficeLamp struct {
	name        string `default:"OfficeLamp"`
	pub         string `default:"set/officelamp"`
	sub         string `default:"reply/officelamp"`
	uniqueConst string `default:"ol"`
}

func (officeLamp *OfficeLamp) Name() string {
	return officeLamp.name
}

func (officeLamp *OfficeLamp) PubTopic() string {
	return officeLamp.pub
}

func (officeLamp *OfficeLamp) SubTopic() string {
	return officeLamp.sub
}

func (officeLamp *OfficeLamp) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(officeLamp.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, officeLamp.Name()+": "+msg)
			if err != nil {
				return
			}

		}()
	}
	return handler, officeLamp.SubTopic()
}

func (officeLamp *OfficeLamp) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[officeLampWhite] = &tb.Btn{Unique: officeLampWhite + officeLamp.uniqueConst, Text: "⬜"}
	buttons[officeLampYellow] = &tb.Btn{Unique: officeLampYellow + officeLamp.uniqueConst, Text: "\U0001F7E8"}
	buttons[officeLampBlue] = &tb.Btn{Unique: officeLampBlue + officeLamp.uniqueConst, Text: "\U0001F7E6"}
	buttons[officeLampGreen] = &tb.Btn{Unique: officeLampGreen + officeLamp.uniqueConst, Text: "\U0001F7E9"}
	buttons[officeLampRed] = &tb.Btn{Unique: officeLampRed + officeLamp.uniqueConst, Text: "\U0001F7E5"}
	buttons[officeLampPink] = &tb.Btn{Unique: officeLampPink + officeLamp.uniqueConst, Text: "\U0001F7EA"}
	buttons[officeLampOff] = &tb.Btn{Unique: officeLampOff + officeLamp.uniqueConst, Text: "🚫"}

	return buttons
}

func (officeLamp *OfficeLamp) Keyboard(services *ServiceContainer) {

	buttons := officeLamp.GenerateKboardBtns()

	officeLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(*buttons[officeLampWhite],
			*buttons[officeLampYellow], *buttons[officeLampBlue],
			*buttons[officeLampGreen], *buttons[officeLampRed],
			*buttons[officeLampPink], *buttons[officeLampOff]),
	)

	officeLamp.AwakenButtons(buttons, services)

	services.botHandler.keyboards[OfficeLampKeyboard] = officeLampModesKeyboard
}

func (officeLamp *OfficeLamp) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(officeLamp.PubTopic(), color)
			})
		}(btn, color)
	}
}

func (officeLamp *OfficeLamp) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officelamp", "Office lamp", OfficeLampKeyboard, KBOARD)
}
