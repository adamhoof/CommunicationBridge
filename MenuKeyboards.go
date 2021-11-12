package main

import tb "gopkg.in/tucnak/telebot.v2"

const (
	ALL_TOYS_KBOARD    = "allToys"
	OFFICE_TOYS_KBOARD = "officeToys"
)

type MenuKeyboards struct {
	keyboards map[string]*tb.ReplyMarkup
}

func (menuKeyboards *MenuKeyboards) OfficeToys(telegramBot *TelegramBot) {
	officeToysKboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[OFFICE_TOYS_KBOARD] = officeToysKboard

	officeLampBtn := officeToysKboard.Text("Table Lamp")
	officeCeilLightBtn := officeToysKboard.Text("Office Ceil Light")
	backBtn := officeToysKboard.Text("⬅ Back")

	officeToysKboard.Reply(
		officeToysKboard.Row(officeLampBtn, officeCeilLightBtn),
		officeToysKboard.Row(backBtn),
	)

	telegramBot.UserEvent(&officeLampBtn, "Office lamp modes", OFFICE_LAMP_KEYBOARD, KBOARD)
	telegramBot.UserEvent(&officeCeilLightBtn, "Office Ceil Light Modes", OFFICE_CEIL_LIGHT_KBOARD, KBOARD)
	telegramBot.UserEvent(&backBtn, "All toys", ALL_TOYS_KBOARD, KBOARD)
}

func (menuKeyboards *MenuKeyboards) AllToys(botHandler *TelegramBot) {
	allToysKboard := &tb.ReplyMarkup{}
	botHandler.keyboards[ALL_TOYS_KBOARD] = allToysKboard

	officeToyBtn := allToysKboard.Text("Office toys")
	bedroomToyBtn := allToysKboard.Text("Bedroom toys")
	cryptoQueryBtn := allToysKboard.Text("Crypto query")

	allToysKboard.Reply(
		allToysKboard.Row(officeToyBtn, bedroomToyBtn),
		allToysKboard.Row(cryptoQueryBtn),
	)

	botHandler.UserEvent("/toys", "/toys", ALL_TOYS_KBOARD, KBOARD)
	botHandler.UserEvent(&officeToyBtn, "Office Toys", OFFICE_TOYS_KBOARD, KBOARD)
	botHandler.UserEvent(&cryptoQueryBtn, "Crypto Query", CRYPTO_QUERY_KBOARD, KBOARD)
}
