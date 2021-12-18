package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	white                = "w"
	yellow               = "y"
	blue                 = "b"
	green                = "g"
	red                  = "r"
	pink                 = "p"
	off                  = "o"
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

	buttons[white] = &tb.Btn{Unique: white, Text: "⬜"}
	buttons[yellow] = &tb.Btn{Unique: yellow, Text: "\U0001F7E8"}
	buttons[blue] = &tb.Btn{Unique: blue, Text: "\U0001F7E6"}
	buttons[green] = &tb.Btn{Unique: green, Text: "\U0001F7E9"}
	buttons[red] = &tb.Btn{Unique: red, Text: "\U0001F7E5"}
	buttons[pink] = &tb.Btn{Unique: pink, Text: "\U0001F7EA"}
	buttons[off] = &tb.Btn{Unique: off, Text: "🚫"}

	return buttons
}

func (officeLamp *OfficeLamp) Kboard(services *ServiceContainer) {

	buttons := officeLamp.GenerateKboardBtns()

	officeLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(*buttons[white],
			*buttons[yellow], *buttons[blue],
			*buttons[green], *buttons[red],
			*buttons[pink], *buttons[off]),
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
