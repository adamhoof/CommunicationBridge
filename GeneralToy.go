package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

type GeneralToy struct {
	name             string
	availableModes   []string
	lastKnownCommand string
	id               int
	publishTopic     string
	subscribeTopic   string
}

func (toy *GeneralToy) Name() string {
	return toy.name
}

func (toy *GeneralToy) PubTopic() string {
	return toy.publishTopic
}

func (toy *GeneralToy) SubTopic() string {
	return toy.subscribeTopic
}

func (toy *GeneralToy) MQTTCommandHandler(services *ServiceContainer) {

	handler := func(client mqtt.Client, message mqtt.Message) {

		func() {
			services.db.UpdateToyMode(toy.Name(), toy.lastKnownCommand)
			_, err := services.botHandler.bot.Send(&me, toy.Name()+": "+toy.lastKnownCommand)
			if err != nil {
				return
			}

		}()
	}
	services.mqtt.SetSubscription(handler, toy.SubTopic())
}

func (toy *GeneralToy) GenerateButtons() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	for _, command := range toy.availableModes {
		func() {
			buttons[command] = &tb.Btn{Unique: command + strconv.Itoa(toy.id), Text: commandAndColorTemplatesForButton[command]}
		}()
	}

	return buttons
}

func (toy *GeneralToy) Keyboard(services *ServiceContainer) {

	buttons := toy.GenerateButtons()
	var buttonsSlice = make([]tb.Btn, len(buttons))

	i := 0
	for name, _ := range buttons {
		buttonsSlice[i] = *buttons[name]
		i++
	}

	keyboard := &tb.ReplyMarkup{ResizeKeyboard: true}
	keyboard.Inline(
		keyboard.Row(buttonsSlice...))

	toy.AwakenButtons(buttons, services)

	services.botHandler.keyboards[toy.Name()] = keyboard
}

func (toy *GeneralToy) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for command, btn := range buttons {

		func(btn *tb.Btn, command string) {

			services.botHandler.bot.Handle(btn, func(c tb.Context) error {
				err := services.botHandler.bot.Respond(c.Callback(), &tb.CallbackResponse{})
				if err != nil {
					return err
				}
				toy.lastKnownCommand = command
				services.mqtt.PublishText(toy.PubTopic(), command)
				return nil
			})
		}(btn, command)
	}
}
