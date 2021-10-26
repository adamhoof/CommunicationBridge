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
	SetKeyboardActions(botHandler *TelegramBotHandler, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateKeyboard(botHandler *TelegramBotHandler) map[string]*tb.Btn
	KeyboardRequestHandler(botHandler *TelegramBotHandler)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, botHandler *TelegramBotHandler,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()

	keyboard := applianceInteractionHandler.GenerateKeyboard(botHandler)
	applianceInteractionHandler.SetKeyboardActions(botHandler, mqttHandler, keyboard)
	applianceInteractionHandler.KeyboardRequestHandler(botHandler)
}

func (botHandler *TelegramBotHandler) OfficeAppliancesKeyboardHandler() {
	officeAppliancesKeyboard := &tb.ReplyMarkup{}
	botHandler.keyboards[OFFICE_APPLIANCES_KEYBOARD] = officeAppliancesKeyboard

	tableLampBtn := officeAppliancesKeyboard.Text(TABLE_LAMP_COMMAND)
	backBtn := officeAppliancesKeyboard.Text("⬅ Back")
	officeAppliancesKeyboard.Reply(
		officeAppliancesKeyboard.Row(tableLampBtn),
		officeAppliancesKeyboard.Row(backBtn),
	)

	botHandler.UserEvent(OFFICE_APPLIANCES_COMMAND, "Office Appliances", OFFICE_APPLIANCES_KEYBOARD, KBOARD)
	botHandler.UserEvent(&tableLampBtn, "Table lamp modes", TABLE_LAMP_KEYBOARD, KBOARD)
	botHandler.UserEvent(&backBtn, "/appliances", ALL_APPLIANCES_KEYBOARD, KBOARD)
}

func (botHandler *TelegramBotHandler) AllAppliancesKeyboardHandler() {
	allAppliancesKeyboard := &tb.ReplyMarkup{}
	botHandler.keyboards[ALL_APPLIANCES_KEYBOARD] = allAppliancesKeyboard

	officeAppliancesBtn := allAppliancesKeyboard.Text("Office appliances")
	bedRoomAppliancesBtn := allAppliancesKeyboard.Text("Bedroom appliances")

	allAppliancesKeyboard.Reply(
		allAppliancesKeyboard.Row(officeAppliancesBtn, bedRoomAppliancesBtn),
	)

	botHandler.UserEvent(&officeAppliancesBtn, OFFICE_APPLIANCES_COMMAND, OFFICE_APPLIANCES_KEYBOARD, KBOARD)
	botHandler.UserEvent(APPLIANCES_COMMAND, APPLIANCES_COMMAND, ALL_APPLIANCES_KEYBOARD, KBOARD)
}
