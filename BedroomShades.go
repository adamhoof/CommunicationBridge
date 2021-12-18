package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BedroomShades struct {
}

const (
	bedroomShadesSub        = "bedroomshades/espReply"
	bedroomShadesPub        = "bedroomshades/rpiSet"
	fullyClose              = "0"
	fullyOpen               = "1"
	BEDROOM_SHADES_KEYBOARD = "bedroomShades"
)

func (bedroomShades *BedroomShades) Name() string {
	return "bedroomShades"
}

func (bedroomShades *BedroomShades) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(bedroomShades.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, msg)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	return handler, bedroomShadesSub
}

func (bedroomShades *BedroomShades) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[fullyOpen] = &tb.Btn{Unique: fullyOpen, Text: "🌞"}
	buttons[fullyClose] = &tb.Btn{Unique: fullyClose, Text: "🌚"}

	return buttons
}

func (bedroomShades *BedroomShades) Kboard(services *ServiceContainer) {
	buttons := bedroomShades.GenerateKboardBtns()

	bedroomShadesModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	bedroomShadesModesKeyboard.Inline(
		bedroomShadesModesKeyboard.Row(*buttons[fullyOpen], *buttons[fullyClose]))

	bedroomShades.AwakenButtons(buttons, services)

	services.botHandler.keyboards[BEDROOM_SHADES_KEYBOARD] = bedroomShadesModesKeyboard
}

func (bedroomShades *BedroomShades) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for mode, btn := range buttons {

		func(btn *tb.Btn, mode string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(bedroomShadesPub, mode)
			})
		}(btn, mode)
	}
}

func (bedroomShades *BedroomShades) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/bedroomshades", "Bedroom shades", BEDROOM_SHADES_KEYBOARD, KBOARD)
}
