package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const TABLE_LAMP_KEYBOARD = "tableLamp"

const (
	white        = "w"
	yellow       = "y"
	blue         = "b"
	green        = "g"
	red          = "r"
	pink         = "p"
	off          = "o"
	tableLampPub = "room/tableLamp/rpiSet"
	tableLampSub = "room/tableLamp/espReply"
)

type OfficeTableLamp struct {
}

func (officeTableLamp *OfficeTableLamp) Name() string {
	return "tableLamp"
}

func (officeTableLamp *OfficeTableLamp) GenerateFunctionButtons(services *ServiceContainer) map[string]*tb.Btn {

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
				services.mqtt.PublishUpdate(tableLampPub, color)
			})
		}(btn, color)
	}
	return buttons
}

func (officeTableLamp *OfficeTableLamp) KeyboardCommands(services *ServiceContainer) {

	buttons := officeTableLamp.GenerateFunctionButtons(services)

	tableLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	tableLampModesKeyboard.Inline(
		tableLampModesKeyboard.Row(*buttons[white],
			*buttons[yellow], *buttons[blue],
			*buttons[green], *buttons[red],
			*buttons[pink], *buttons[off]),
	)
	services.botHandler.keyboards[TABLE_LAMP_KEYBOARD] = tableLampModesKeyboard
}

func (officeTableLamp *OfficeTableLamp) MQTTMessageProcessor(services *ServiceContainer) (TableLampMessageHandler mqtt.MessageHandler, topic string) {

	TableLampMessageHandler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateMode("TableLamp", string(message.Payload()))
		}()
	}
	return TableLampMessageHandler, tableLampSub
}

func (officeTableLamp *OfficeTableLamp) NonKeyboardCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/tablelamp", "Table Lamp", TABLE_LAMP_KEYBOARD, KBOARD)
}
