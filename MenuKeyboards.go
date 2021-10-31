package main

import tb "gopkg.in/tucnak/telebot.v2"

const (
	ALL_APPLIANCES_KEYBOARD    = "allAppliances"
	OFFICE_APPLIANCES_KEYBOARD = "officeAppliances"
	CRYPTO_DATA_KEYBOARD       = "cryptoData"
)

type MenuKeyboards struct {
	keyboards map[string]*tb.ReplyMarkup
}

func (menuKeyboards *MenuKeyboards) OfficeAppliances(telegramBot *TelegramBot) {
	officeAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[OFFICE_APPLIANCES_KEYBOARD] = officeAppliancesKeyboard

	tableLampBtn := officeAppliancesKeyboard.Text("Table Lamp")
	backBtn := officeAppliancesKeyboard.Text("⬅ Back")
	officeAppliancesKeyboard.Reply(
		officeAppliancesKeyboard.Row(tableLampBtn),
		officeAppliancesKeyboard.Row(backBtn),
	)

	telegramBot.UserEvent(&tableLampBtn, "Table lamp modes", TABLE_LAMP_KEYBOARD, KBOARD)
	telegramBot.UserEvent(&backBtn, "Appliances", ALL_APPLIANCES_KEYBOARD, KBOARD)
}

func (menuKeyboards *MenuKeyboards) AllToys(botHandler *TelegramBot) {
	allAppliancesKeyboard := &tb.ReplyMarkup{}
	botHandler.keyboards[ALL_APPLIANCES_KEYBOARD] = allAppliancesKeyboard

	officeAppliancesBtn := allAppliancesKeyboard.Text("Office appliances")
	bedRoomAppliancesBtn := allAppliancesKeyboard.Text("Bedroom appliances")
	cryptoDataQueryBtn := allAppliancesKeyboard.Text("Crypto data")

	allAppliancesKeyboard.Reply(
		allAppliancesKeyboard.Row(officeAppliancesBtn, bedRoomAppliancesBtn),
		allAppliancesKeyboard.Row( cryptoDataQueryBtn),
	)

	botHandler.UserEvent("/appliances", "/appliances", ALL_APPLIANCES_KEYBOARD, KBOARD)
	botHandler.UserEvent(&officeAppliancesBtn, "Office Appliances", OFFICE_APPLIANCES_KEYBOARD, KBOARD)
	botHandler.UserEvent(&cryptoDataQueryBtn, "Crypto Query", CRYPTO_DATA_KEYBOARD, KBOARD)
}
