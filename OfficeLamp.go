package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const OFFICE_LAMP_KEYBOARD = "officelamp"

const (
	white        = "w"
	yellow       = "y"
	blue         = "b"
	green        = "g"
	red          = "r"
	pink         = "p"
	off          = "o"
	officeLampPub = "room/officeLamp/rpiSet"
	officeLampSub = "room/officeLamp/espReply"
)

type OfficeLamp struct {
}

func (officeLamp *OfficeLamp) Name() string {
	return "officeLamp"
}

func (officeLamp *OfficeLamp) CreateDBObject(toyName string, toyMode string, postgreHandler *PostgreSQLHandler) {
	postgreHandler.CreateToy(toyName, toyMode)
}

func (officeLamp *OfficeLamp) GenerateFunctionButtons(services *ServiceContainer) map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[white] = &tb.Btn{Unique: white, Text: "⬜"}
	buttons[yellow] = &tb.Btn{Unique: yellow, Text: "\U0001F7E8"}
	buttons[blue] = &tb.Btn{Unique: blue, Text: "\U0001F7E6"}
	buttons[green] = &tb.Btn{Unique: green, Text: "\U0001F7E9"}
	buttons[red] = &tb.Btn{Unique: red, Text: "\U0001F7E5"}
	buttons[pink] = &tb.Btn{Unique: pink, Text: "\U0001F7EA"}
	buttons[off] = &tb.Btn{Unique: off, Text: "🚫"}

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
	return buttons
}

func (officeLamp *OfficeLamp) KeyboardCommands(services *ServiceContainer) {

	buttons := officeLamp.GenerateFunctionButtons(services)

	officeLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeLampModesKeyboard.Inline(
		officeLampModesKeyboard.Row(*buttons[white],
			*buttons[yellow], *buttons[blue],
			*buttons[green], *buttons[red],
			*buttons[pink], *buttons[off]),
	)
	services.botHandler.keyboards[OFFICE_LAMP_KEYBOARD] = officeLampModesKeyboard
}

func (officeLamp *OfficeLamp) MQTTMessageProcessor(services *ServiceContainer) (OfficeLampMessageHandler mqtt.MessageHandler, topic string) {

	OfficeLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateToyMode(officeLamp.Name(), string(message.Payload()))
		}()
	}
	return OfficeLampMessageHandler, officeLampSub
}

func (officeLamp *OfficeLamp) NonKeyboardCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officelamp", "Office lamp", OFFICE_LAMP_KEYBOARD, KBOARD)
}
