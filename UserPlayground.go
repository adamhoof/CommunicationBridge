package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

var toyColors = make(map[string]string)

type Playground struct {
}

type ToyAttributes struct {
	name            string
	commandWithName []string
	id              int
	publishTopic    string
	subscribeTopic  string
}

func (playground *Playground) ColorTheToys() {
	toyColors["on"] = "⬜"
	toyColors["white"] = "⬜"
	toyColors["yellow"] = "\U0001F7E8"
	toyColors["blue"] = "\U0001F7E6"
	toyColors["green"] = "\U0001F7E9"
	toyColors["red"] = "\U0001F7E5"
	toyColors["pink"] = "\U0001F7EA"
	toyColors["orange"] = "\U0001F7E7"
	toyColors["off"] = "🚫"
}

func (toyAttributes *ToyAttributes) Name() string {
	return toyAttributes.name
}

func (toyAttributes *ToyAttributes) PubTopic() string {
	return toyAttributes.publishTopic
}

func (toyAttributes *ToyAttributes) SubTopic() string {
	return toyAttributes.subscribeTopic
}

func (toyAttributes *ToyAttributes) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			msg := string(message.Payload())
			services.db.UpdateToyMode(toyAttributes.Name(), msg)
			_, err := services.botHandler.bot.Send(&me, toyAttributes.Name()+": "+msg)
			if err != nil {
				return
			}

		}()
	}
	return handler, toyAttributes.SubTopic()
}

func (toyAttributes *ToyAttributes) GenerateButtons() map[string]*tb.Btn {

	buttons := make(map[string]*tb.Btn)

	for _, command := range toyAttributes.commandWithName {
		func() {
			buttons[command] = &tb.Btn{Unique: command + strconv.Itoa(toyAttributes.id), Text: toyColors[command]}
		}()
	}

	return buttons
}

func (toyAttributes *ToyAttributes) Keyboard(services *ServiceContainer) {

	buttons := toyAttributes.GenerateButtons()
	var buttonsSlice = make([]tb.Btn, len(buttons))

	i := 0
	for name, _ := range buttons {
		buttonsSlice[i] = *buttons[name]
		i++
	}

	keyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	keyboard.Inline(
		keyboard.Row(buttonsSlice...))

	toyAttributes.AwakenButtons(buttons, services)

	services.botHandler.keyboards[toyAttributes.Name()] = keyboard
}

func (toyAttributes *ToyAttributes) AwakenButtons(buttons map[string]*tb.Btn, services *ServiceContainer) {

	for mode, btn := range buttons {

		func(btn *tb.Btn, mode string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}
				services.mqtt.PublishText(toyAttributes.PubTopic(), mode)
			})
		}(btn, mode)
	}
}

func (toyAttributes *ToyAttributes) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/officelamp", "Office lamp", OfficeLampKeyboard, KBOARD)
}

type Toy interface {
	Name() string
	PubTopic() string
	SubTopic() string
	MQTTCommandHandler(services *ServiceContainer) (MessageHandler mqtt.MessageHandler, topic string)
	Keyboard(services *ServiceContainer)
	TextCommands(services *ServiceContainer)
}

func (playground *Playground) spawnToys(services *ServiceContainer) {

}

func (playground *Playground) takeOutToys(toyStorage *ToyBag, services *ServiceContainer) {

	for _, toy := range toyStorage.bag {
		handler, topic := toy.MQTTCommandHandler(services)
		services.mqtt.SetSubscription(handler, topic)

		toy.Keyboard(services)
		toy.TextCommands(services)

		services.db.CreateToy(toy.Name(), "")
	}
}
