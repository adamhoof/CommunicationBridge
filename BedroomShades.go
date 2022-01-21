package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	bedroomShadesFullyOpen  = "1"
	bedroomShadesFullyClose = "0"
	BEDROOM_SHADES_KEYBOARD = "bedroomShades"
)

type BedroomShades struct {
	name        string `default:"BedroomShades"`
	pub         string `default:"set/bedroomshades"`
	sub         string `default:"reply/bedroomshades"`
	uniqueConst string `default:"bs"`
}

func (bedroomShades *BedroomShades) Name() string {
	return bedroomShades.name
}

func (bedroomShades *BedroomShades) PubTopic() string {
	return bedroomShades.pub
}

func (bedroomShades *BedroomShades) SubTopic() string {
	return bedroomShades.sub
}
func (bedroomShades *BedroomShades) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			if msg != "done" {
				services.db.UpdateToyMode(bedroomShades.Name(), msg)
			}
			_, err := services.botHandler.bot.Send(&me, bedroomShades.Name()+": "+msg)
			if err != nil {
				return
			}
		}()
	}
	return handler, bedroomShades.SubTopic()
}

func (bedroomShades *BedroomShades) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[bedroomShadesFullyOpen] = &tb.Btn{Unique: bedroomShadesFullyOpen + "bsb", Text: "🌞"}
	buttons[bedroomShadesFullyClose] = &tb.Btn{Unique: bedroomShadesFullyClose + "bsb", Text: "🌚"}

	return buttons
}

func (bedroomShades *BedroomShades) Keyboard(services *ServiceContainer) {
	buttons := bedroomShades.GenerateKboardBtns()

	bedroomShadesModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	bedroomShadesModesKeyboard.Inline(
		bedroomShadesModesKeyboard.Row(*buttons[bedroomShadesFullyOpen], *buttons[bedroomShadesFullyClose]))

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
				services.mqtt.PublishText(bedroomShades.PubTopic(), mode)
			})
		}(btn, mode)
	}
}

func (bedroomShades *BedroomShades) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/bedroomshades", "Bedroom shades", BEDROOM_SHADES_KEYBOARD, KBOARD)
}
