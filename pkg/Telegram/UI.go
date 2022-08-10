package telegram

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	tb "gopkg.in/telebot.v3"
)

const (
	AllToysKeyboard     = "allToys"
	OfficeToysKeyboard  = "officeToys"
	BedroomToysKeyboard = "bedroomToys"
)

func CreateRoomOfToysKeyboard(handler *BotHandler, keyboards map[string]*tb.ReplyMarkup, toys map[string]*connectable.Toy, name string, buttonNameUnificator string) {
	keyboard := &tb.ReplyMarkup{}
	keyboards[name] = keyboard

	var buttons []tb.Btn
	for _, toy := range toys {
		button := tb.Btn{Text: buttonNameUnificator + toy.Name}
		buttons = append(buttons, button)
	}
	backBtn := keyboard.Text("⬅")
	handler.SendKeyboardOnButtonClick(&backBtn, "⬅", keyboards, AllToysKeyboard)

	keyboard.Reply(
		keyboard.Row(buttons...),
		keyboard.Row(backBtn))
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
