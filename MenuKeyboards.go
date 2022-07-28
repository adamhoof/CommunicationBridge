package main

/*




import (
	tb "gopkg.in/telebot.v3"
)

const (
	AllToysKeyboard     = "allToys"
	OfficeToysKeyboard  = "officeToys"
	BedroomToysKeyboard = "bedroomToys"
)

type MenuKeyboards struct {
	keyboards map[string]*tb.ReplyMarkup
}

func (menuKeyboards *MenuKeyboards) OfficeToys(telegramBot *TelegramBot) {
	officeToysKeyboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[OfficeToysKeyboard] = officeToysKeyboard

	officeLampBtn := officeToysKeyboard.Text("Table Lamp")
	officeCeilLightBtn := officeToysKeyboard.Text("Main Light")
	backBtn := officeToysKeyboard.Text("⬅")

	officeToysKeyboard.Reply(
		officeToysKeyboard.Row(officeLampBtn, officeCeilLightBtn),
		officeToysKeyboard.Row(backBtn),
	)

	telegramBot.UserEvent(&officeLampBtn, "Office lamp modes", "OfficeLamp", KBOARD)
	telegramBot.UserEvent(&officeCeilLightBtn, "Office Ceil Light Modes", "OfficeCeilLight", KBOARD)
	telegramBot.UserEvent(&backBtn, "All toys", AllToysKeyboard, KBOARD)
}

func (menuKeyboards *MenuKeyboards) BedroomToys(telegramBot *TelegramBot) {
	bedroomToysKboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[BedroomToysKeyboard] = bedroomToysKboard

	bedroomShadesBtn := bedroomToysKboard.Text("Shades")
	bedroomLampBtn := bedroomToysKboard.Text("Lamp")
	backBtn := bedroomToysKboard.Text("⬅")

	bedroomToysKboard.Reply(
		bedroomToysKboard.Row(bedroomShadesBtn, bedroomLampBtn),
		bedroomToysKboard.Row(backBtn))

	telegramBot.UserEvent(&bedroomShadesBtn, "Bedroom shades modes", "BedroomShades", KBOARD)
	telegramBot.UserEvent(&bedroomLampBtn, "Bedroom lamp modes", "BedroomLamp", KBOARD)
	telegramBot.UserEvent(&backBtn, "All toys", AllToysKeyboard, KBOARD)
}

func (menuKeyboards *MenuKeyboards) AllToys(botHandler *TelegramBot) {
	allToysKboard := &tb.ReplyMarkup{}
	botHandler.keyboards[AllToysKeyboard] = allToysKboard

	officeToyBtn := allToysKboard.Text("Office toys")
	bedroomToyBtn := allToysKboard.Text("Bedroom toys")

	allToysKboard.Reply(
		allToysKboard.Row(officeToyBtn, bedroomToyBtn),
	)

	botHandler.UserEvent("/toys", "/toys", AllToysKeyboard, KBOARD)
	botHandler.UserEvent(&officeToyBtn, "Office Toys", OfficeToysKeyboard, KBOARD)
	botHandler.UserEvent(&bedroomToyBtn, "Bedroom Toys", BedroomToysKeyboard, KBOARD)
}
*/
