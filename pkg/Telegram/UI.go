package telegram

import (
	tb "gopkg.in/telebot.v3"
)

const (
	AllToysKeyboard     = "allToys"
	OfficeToysKeyboard  = "officeToys"
	BedroomToysKeyboard = "bedroomToys"
)

func CreateOfficeToysKeyboardUI(handler *BotHandler, keyboards map[string]*tb.ReplyMarkup) {
	officeToysKeyboard := &tb.ReplyMarkup{}
	keyboards[OfficeToysKeyboard] = officeToysKeyboard

	officeLampBtn := officeToysKeyboard.Text("o Table Lamp")
	officeCeilLightBtn := officeToysKeyboard.Text("o Ceil Light")
	backBtn := officeToysKeyboard.Text("⬅")

	officeToysKeyboard.Reply(
		officeToysKeyboard.Row(officeLampBtn, officeCeilLightBtn),
		officeToysKeyboard.Row(backBtn),
	)

	handler.SendKeyboardOnButtonClick(&officeLampBtn, "o Table Lamp modes", keyboards["OfficeLamp"])
	handler.SendKeyboardOnButtonClick(&officeCeilLightBtn, "o Office Light modes", keyboards["OfficeCeilLight"])
	handler.SendKeyboardOnButtonClick(&backBtn, "⬅", keyboards[AllToysKeyboard])
}

func CreateBedroomToysKeyboardUI(handler *BotHandler, keyboards map[string]*tb.ReplyMarkup) {
	bedroomToysKeyboard := &tb.ReplyMarkup{}
	keyboards["BedroomToysKeyboard"] = bedroomToysKeyboard

	bedroomShadesBtn := bedroomToysKeyboard.Text("b Shades")
	bedroomLampBtn := bedroomToysKeyboard.Text("b Table Lamp")
	backBtn := bedroomToysKeyboard.Text("⬅")

	bedroomToysKeyboard.Reply(
		bedroomToysKeyboard.Row(bedroomShadesBtn, bedroomLampBtn),
		bedroomToysKeyboard.Row(backBtn))

	handler.SendKeyboardOnButtonClick(&bedroomShadesBtn, "b Shades modes", keyboards["BedroomShades"])
	handler.SendKeyboardOnButtonClick(&bedroomLampBtn, "b Table Lamp modes", keyboards["BedroomLamp"])
	handler.SendKeyboardOnButtonClick(&backBtn, "⬅", keyboards[AllToysKeyboard])
}

func CreateAllToysKeyboardUI(handler *BotHandler, keyboards map[string]*tb.ReplyMarkup) {
	allToysKeyboard := &tb.ReplyMarkup{}
	keyboards[AllToysKeyboard] = allToysKeyboard

	officeToysBtn := allToysKeyboard.Text("Office toys")
	bedroomToysBtn := allToysKeyboard.Text("Bedroom toys")

	allToysKeyboard.Reply(
		allToysKeyboard.Row(officeToysBtn, bedroomToysBtn),
	)

	handler.SendKeyboardOnButtonClick(&officeToysBtn, "Office Toys", keyboards[OfficeToysKeyboard])
	handler.SendKeyboardOnButtonClick(&bedroomToysBtn, "Bedroom Toys", keyboards[BedroomToysKeyboard])
}
