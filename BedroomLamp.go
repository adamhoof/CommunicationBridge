package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	bedroomLampYellow     = "yellow"
	bedroomLampOrange     = "orange"
	bedroomLampBlue       = "blue"
	bedroomLampGreen      = "green"
	bedroomLampRed        = "red"
	bedroomLampPink       = "pink"
	bedroomLampOff        = "off"
	bedroomLampPub        = "bedroomlamp/rpiSet"
	bedroomLampSub        = "bedroomlamp/espReply"
	BEDROOM_LAMP_KEYBOARD = "bedroomlamp"
)

type BedroomLamp struct {
}

func (bedroomLamp *BedroomLamp) Name() string {
	return "Bedroom Lamp"
}

func (bedroomLamp *BedroomLamp) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(bedroomLamp.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, bedroomLamp.Name()+": "+msg)
			if err != nil {
				return
			}

		}()
	}
	return handler, bedroomLampSub
}

func (bedroomLamp *BedroomLamp) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[bedroomLampYellow] = &tb.Btn{Unique: bedroomLampYellow + "bl", Text: "\U0001F7E8"}
	buttons[bedroomLampOrange] = &tb.Btn{Unique: bedroomLampOrange + "bl", Text: "\U0001F7E7"}
	buttons[bedroomLampBlue] = &tb.Btn{Unique: bedroomLampBlue + "bl", Text: "\U0001F7E6"}
	buttons[bedroomLampGreen] = &tb.Btn{Unique: bedroomLampGreen + "bl", Text: "\U0001F7E9"}
	buttons[bedroomLampRed] = &tb.Btn{Unique: bedroomLampRed + "bl", Text: "\U0001F7E5"}
	buttons[bedroomLampPink] = &tb.Btn{Unique: bedroomLampPink + "bl", Text: "\U0001F7EA"}
	buttons[bedroomLampOff] = &tb.Btn{Unique: bedroomLampOff + "bl", Text: "🚫"}

	return buttons
}

func (bedroomLamp *BedroomLamp) Kboard(services *ServiceContainer) {

	buttons := bedroomLamp.GenerateKboardBtns()

	bedroomLampModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	bedroomLampModesKeyboard.Inline(
		bedroomLampModesKeyboard.Row(
			*buttons[bedroomLampYellow], *buttons[bedroomLampOrange], *buttons[bedroomLampBlue],
			*buttons[bedroomLampGreen], *buttons[bedroomLampRed],
			*buttons[bedroomLampPink], *buttons[bedroomLampOff]),
	)

	bedroomLamp.AwakenButtons(buttons, services)

	services.botHandler.keyboards[BEDROOM_LAMP_KEYBOARD] = bedroomLampModesKeyboard
}

func (bedroomLamp *BedroomLamp) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(bedroomLampPub, color)
			})
		}(btn, color)
	}
}

func (bedroomLamp *BedroomLamp) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/bedroomlamp", "Bedroom lamp", BEDROOM_LAMP_KEYBOARD, KBOARD)
}
