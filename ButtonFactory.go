package main

import (
	"errors"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

type ButtonFactory struct {
	commandAndStyleForButtonTemplates map[string]string
}

func (factory *ButtonFactory) setCommandAndIconButtonTemplates(commandAndColorTemplates map[string]string) {
	factory.commandAndStyleForButtonTemplates = commandAndColorTemplates
}

func (factory *ButtonFactory) generateInlineButtons(unificator int, commands []string) (map[string]*tb.Btn, error) {

	if factory.commandAndStyleForButtonTemplates == nil {
		err := errors.New("no command and style template provided, use function setCommandAndIconButtonTemplates")
		return nil, err
	}

	buttons := make(map[string]*tb.Btn)

	for _, command := range commands {
		func() {
			buttons[command] = &tb.Btn{Unique: command + strconv.Itoa(unificator), Text: factory.commandAndStyleForButtonTemplates[command]}
		}()
	}

	return buttons, nil
}

func (factory *ButtonFactory) setButtonHandlers(buttons map[string]*tb.Btn, toy *Toy, services *ServiceContainer) {

	for command, btn := range buttons {

		func(btn *tb.Btn, command string) {

			services.botHandler.bot.Handle(btn, func(c tb.Context) error {
				err := services.botHandler.bot.Respond(c.Callback(), &tb.CallbackResponse{})
				if err != nil {
					return err
				}
				toy.lastKnownCommand = command
				services.mqtt.PublishText(toy.publishTopic, command)
				return nil
			})
		}(btn, command)
	}
}
