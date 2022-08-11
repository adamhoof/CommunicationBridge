package telegram

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

const (
	AllToysKeyboard     = "allToys"
	OfficeToysKeyboard  = "officeToys"
	BedroomToysKeyboard = "bedroomToys"
)

var toyCommandIconPairs = map[string]string{
	"on":     "⬜",
	"white":  "⬜",
	"yellow": "\U0001F7E8",
	"blue":   "\U0001F7E6",
	"green":  "\U0001F7E9",
	"red":    "\U0001F7E5",
	"pink":   "\U0001F7EA",
	"orange": "\U0001F7E7",
	"off":    "🚫",
	"open":   "🌞",
	"close":  "🌚"}

func CreateReplyButtonsForToysInRoom(toysGroup map[string]*connectable.Toy, buttonNameUnificator string) (buttons []tb.Btn) {
	for _, toy := range toysGroup {
		button := tb.Btn{Text: buttonNameUnificator + toy.Name}
		buttons = append(buttons, button)
	}
	return buttons
}

func CreateReplyKeyboardFromButtons(buttons []tb.Btn, backButton tb.Btn) (keyboard *tb.ReplyMarkup) {
	keyboard.Reply(
		keyboard.Row(buttons...),
		keyboard.Row(backButton),
	)
	return keyboard
}

func GenerateInlineButtonsForToy(toy *connectable.Toy) map[string]*tb.Btn {
	buttons := make(map[string]*tb.Btn)

	for _, command := range toy.AvailableModes {
		func() {
			buttons[command] = &tb.Btn{Unique: command + strconv.Itoa(toy.Id), Text: toyCommandIconPairs[command]}
		}()
	}
	return buttons
}

func GenerateInlineKeyboardFromButtonsForToy(buttons map[string]*tb.Btn) (keyboard *tb.ReplyMarkup) {
	var buttonsArray = make([]tb.Btn, len(buttons))

	for name, _ := range buttons {
		buttonsArray = append(buttonsArray, *buttons[name])
	}

	keyboard.ResizeKeyboard = true
	keyboard.Inline(
		keyboard.Row(buttonsArray...))

	return keyboard
}

func CreateAllToysKeyboardUI(handler *BotHandler, keyboards map[string]*tb.ReplyMarkup) {
	allToysKeyboard := &tb.ReplyMarkup{}
	keyboards[AllToysKeyboard] = allToysKeyboard

	officeToysBtn := allToysKeyboard.Text("Office toys")
	bedroomToysBtn := allToysKeyboard.Text("Bedroom toys")

	allToysKeyboard.Reply(
		allToysKeyboard.Row(officeToysBtn, bedroomToysBtn),
	)

	handler.SendKeyboardOnButtonClick(&officeToysBtn, "Office Toys", keyboards, OfficeToysKeyboard)
	handler.SendKeyboardOnButtonClick(&bedroomToysBtn, "Bedroom Toys", keyboards, BedroomToysKeyboard)
}
