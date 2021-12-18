package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type OfficeCeilLight struct {
}

const (
	officeCeilLightPub       = "officeceil/rpiSet"
	officeCeilLightSub       = "officeceil/espReply"
	officeCeilLightOn        = "on"
	officeCeilLightOff       = "off"
	OFFICE_CEIL_LIGHT_KBOARD = "officeCeilLight"
)

func (officeCeilLight *OfficeCeilLight) Name() string {
	return "officeceillight"
}

func (officeCeilLight *OfficeCeilLight) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateToyMode(officeCeilLight.Name(), string(message.Payload()))
		}()
	}
	return handler, officeCeilLightSub
}

func (officeCeilLight *OfficeCeilLight) GenerateKboardBtns() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	buttons[officeCeilLightOn] = &tb.Btn{Unique: officeCeilLightOn + "oclb", Text: "⬜"}
	buttons[officeCeilLightOff] = &tb.Btn{Unique: officeCeilLightOff + "oclb", Text: "🚫"}

	return buttons
}

func (officeCeilLight *OfficeCeilLight) Kboard(services *ServiceContainer) {

	buttons := officeCeilLight.GenerateKboardBtns()

	officeCeilLightModesKeyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	officeCeilLightModesKeyboard.Inline(
		officeCeilLightModesKeyboard.Row(*buttons[officeCeilLightOn],
			*buttons[officeCeilLightOff],
		))

	officeCeilLight.AwakenButtons(buttons, services)

	services.botHandler.keyboards[OFFICE_CEIL_LIGHT_KBOARD] = officeCeilLightModesKeyboard
}

func (officeCeilLight *OfficeCeilLight) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for color, btn := range buttons {

		func(btn *tb.Btn, color string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(officeCeilLightPub, color)
			})
		}(btn, color)
	}
}

func (officeCeilLight *OfficeCeilLight) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officeceil", "Office ceil light", OFFICE_CEIL_LIGHT_KBOARD, KBOARD)
}
