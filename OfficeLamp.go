package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	officeLampWhite      = "white"
	officeLampYellow     = "yellow"
	officeLampBlue       = "blue"
	officeLampGreen      = "green"
	officeLampRed        = "red"
	officeLampPink       = "pink"
	officeLampOff        = "off"
	officeLampPub        = "officelamp/rpiSet"
	officeLampSub        = "officelamp/espReply"
	OFFICE_LAMP_KEYBOARD = "officelamp"
)

type OfficeLamp struct {
}

func (officeLamp *OfficeLamp) Name() string {
	return "officelamp"
}

func (officeLamp *OfficeLamp) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateToyMode(officeLamp.Name(), string(message.Payload()))
		}()
	}
	return handler, officeLampSub
}

func (officeLamp *OfficeLamp) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[officeLampWhite] = &tb.Btn{Unique: officeLampWhite + "ol", Text: "⬜"}
	buttons[officeLampYellow] = &tb.Btn{Unique: officeLampYellow + "ol", Text: "\U0001F7E8"}
	buttons[officeLampBlue] = &tb.Btn{Unique: officeLampBlue + "ol", Text: "\U0001F7E6"}
	buttons[officeLampGreen] = &tb.Btn{Unique: officeLampGreen + "ol", Text: "\U0001F7E9"}
	buttons[officeLampRed] = &tb.Btn{Unique: officeLampRed + "ol", Text: "\U0001F7E5"}
	buttons[officeLampPink] = &tb.Btn{Unique: officeLampPink + "ol", Text: "\U0001F7EA"}
	buttons[officeLampOff] = &tb.Btn{Unique: officeLampOff + "ol", Text: "🚫"}

	return buttons
}

func (officeLamp *OfficeLamp) Kboard(services *ServiceContainer) {

	buttons := officeLamp.GenerateKboardBtns()

	officeLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(*buttons[officeLampWhite],
			*buttons[officeLampYellow], *buttons[officeLampBlue],
			*buttons[officeLampGreen], *buttons[officeLampRed],
			*buttons[officeLampPink], *buttons[officeLampOff]),
	)

	officeLamp.AwakenButtons(buttons, services)

	services.botHandler.keyboards[OFFICE_LAMP_KEYBOARD] = officeLampModesKeyboard
}

func (officeLamp *OfficeLamp) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(officeLampPub, color)
			})
		}(btn, color)
	}
}

func (officeLamp *OfficeLamp) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officelamp", "Office lamp", OFFICE_LAMP_KEYBOARD, KBOARD)
}
