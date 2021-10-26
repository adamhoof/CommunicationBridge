package main

import tb "gopkg.in/tucnak/telebot.v2"

type KeyboardsController struct {
}

func (keyboardsController *KeyboardsController) OfficeAppliancesKeyboardHandler(telegramBot *TelegramBot) {
	officeAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[OFFICE_APPLIANCES_KEYBOARD] = officeAppliancesKeyboard

	tableLampBtn := officeAppliancesKeyboard.Text(TABLE_LAMP_COMMAND)
	backBtn := officeAppliancesKeyboard.Text("⬅ Back")
	officeAppliancesKeyboard.Reply(
		officeAppliancesKeyboard.Row(tableLampBtn),
		officeAppliancesKeyboard.Row(backBtn),
	)

	telegramBot.UserEvent(OFFICE_APPLIANCES_COMMAND, "Office Appliances", OFFICE_APPLIANCES_KEYBOARD, KBOARD)
	telegramBot.UserEvent(&tableLampBtn, "Table lamp modes", TABLE_LAMP_KEYBOARD, KBOARD)
	telegramBot.UserEvent(&backBtn, "/appliances", ALL_APPLIANCES_KEYBOARD, KBOARD)
}

func (keyboardsController *KeyboardsController) AllAppliancesKeyboardHandler(telegramBot *TelegramBot) {
	allAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBot.keyboards[ALL_APPLIANCES_KEYBOARD] = allAppliancesKeyboard

	officeAppliancesBtn := allAppliancesKeyboard.Text("Office appliances")
	bedRoomAppliancesBtn := allAppliancesKeyboard.Text("Bedroom appliances")

	allAppliancesKeyboard.Reply(
		allAppliancesKeyboard.Row(officeAppliancesBtn, bedRoomAppliancesBtn),
	)

	telegramBot.UserEvent(&officeAppliancesBtn, OFFICE_APPLIANCES_COMMAND, OFFICE_APPLIANCES_KEYBOARD, KBOARD)
	telegramBot.UserEvent(APPLIANCES_COMMAND, APPLIANCES_COMMAND, ALL_APPLIANCES_KEYBOARD, KBOARD)
}
