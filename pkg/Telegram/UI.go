package telegrambot

import (
	tb "gopkg.in/telebot.v3"
)

const (
	AllToysKeyboard     = "allToys"
	OfficeToysKeyboard  = "officeToys"
	BedroomToysKeyboard = "bedroomToys"
)

func CreateOfficeToysKeyboardUI(handler *BotHandler, keyboardStorage *KeyboardStorage) {
	officeToysKeyboard := &tb.ReplyMarkup{}
	keyboardStorage.AddKeyboard(OfficeToysKeyboard, officeToysKeyboard)

	officeLampBtn := officeToysKeyboard.Text("o Table Lamp")
	officeCeilLightBtn := officeToysKeyboard.Text("o Ceil Light")
	backBtn := officeToysKeyboard.Text("⬅")

	officeToysKeyboard.Reply(
		officeToysKeyboard.Row(officeLampBtn, officeCeilLightBtn),
		officeToysKeyboard.Row(backBtn),
	)

	handler.SendKeyboardOnButtonClick(&officeLampBtn, "o Table Lamp modes", keyboardStorage.GetKeyboardByName("OfficeLamp"))
	handler.SendKeyboardOnButtonClick(&officeCeilLightBtn, "o Office Light modes", keyboardStorage.GetKeyboardByName("OfficeCeilLight"))
	handler.SendKeyboardOnButtonClick(&backBtn, "⬅", keyboardStorage.GetKeyboardByName(AllToysKeyboard))
}

func CreateBedroomToysKeyboardUI(handler *BotHandler, keyboardStorage *KeyboardStorage) {
	bedroomToysKeyboard := &tb.ReplyMarkup{}
	keyboardStorage.AddKeyboard("BedroomToysKeyboard", bedroomToysKeyboard)

	bedroomShadesBtn := bedroomToysKeyboard.Text("b Shades")
	bedroomLampBtn := bedroomToysKeyboard.Text("b Table Lamp")
	backBtn := bedroomToysKeyboard.Text("⬅")

	bedroomToysKeyboard.Reply(
		bedroomToysKeyboard.Row(bedroomShadesBtn, bedroomLampBtn),
		bedroomToysKeyboard.Row(backBtn))

	handler.SendKeyboardOnButtonClick(&bedroomShadesBtn, "b Shades modes", keyboardStorage.GetKeyboardByName("BedroomShades"))
	handler.SendKeyboardOnButtonClick(&bedroomLampBtn, "b Table Lamp modes", keyboardStorage.GetKeyboardByName("BedroomLamp"))
	handler.SendKeyboardOnButtonClick(&backBtn, "⬅", keyboardStorage.GetKeyboardByName(AllToysKeyboard))
}

func CreateAllToysKeyboardUI(handler *BotHandler, keyboardStorage *KeyboardStorage) {
	allToysKeyboard := &tb.ReplyMarkup{}
	keyboardStorage.AddKeyboard(AllToysKeyboard, allToysKeyboard)

	officeToysBtn := allToysKeyboard.Text("Office toys")
	bedroomToysBtn := allToysKeyboard.Text("Bedroom toys")

	allToysKeyboard.Reply(
		allToysKeyboard.Row(officeToysBtn, bedroomToysBtn),
	)

	handler.SendKeyboardOnButtonClick(&officeToysBtn, "Office Toys", keyboardStorage.GetKeyboardByName(OfficeToysKeyboard))
	handler.SendKeyboardOnButtonClick(&bedroomToysBtn, "Bedroom Toys", keyboardStorage.GetKeyboardByName(BedroomToysKeyboard))
}
