package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	APPLIANCES_COMMAND         = "/appliances"
	ALL_APPLIANCES_KEYBOARD    = "allAppliances"
	OFFICE_APPLIANCES_COMMAND  = "/officeappliances"
	OFFICE_APPLIANCES_KEYBOARD = "officeAppliances"
	TABLE_LAMP_COMMAND         = "/tablelamp"
	TABLE_LAMP_KEYBOARD        = "tableLamp"
)

type ApplianceInteractionHandler interface {
	Name() string
	MessageProcessor() (MessageHandler mqtt.MessageHandler)
	SetKeyboardActions(telegramBot *TelegramBot, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateKeyboard(telegramBot *TelegramBot) map[string]*tb.Btn
	KeyboardRequestHandler(telegramBot *TelegramBot)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, telegramBot *TelegramBot,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()

	keyboard := applianceInteractionHandler.GenerateKeyboard(telegramBot)
	applianceInteractionHandler.SetKeyboardActions(telegramBot, mqttHandler, keyboard)
	applianceInteractionHandler.KeyboardRequestHandler(telegramBot)
}

func (telegramBot *TelegramBot) OfficeAppliancesKeyboardHandler() {
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

func (telegramBot *TelegramBot) AllAppliancesKeyboardHandler() {
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
