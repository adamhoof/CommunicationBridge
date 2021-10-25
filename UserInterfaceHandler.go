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
	SetKeyboardActions(telegramBotHandler *TelegramBotHandler, mqttHandler *MQTTHandler, buttons map[string]*tb.Btn)
	GenerateKeyboard(telegramBotHandler *TelegramBotHandler) map[string]*tb.Btn
	KeyboardRequestHandler(botHandler *TelegramBotHandler)
}

func SetupClientInterfaceOptions(applianceInteractionHandler ApplianceInteractionHandler, telegramBotHandler *TelegramBotHandler,
	mqttHandler *MQTTHandler, messageProcessors map[string]mqtt.MessageHandler) {

	messageProcessors[applianceInteractionHandler.Name()] = applianceInteractionHandler.MessageProcessor()

	keyboard := applianceInteractionHandler.GenerateKeyboard(telegramBotHandler)
	applianceInteractionHandler.SetKeyboardActions(telegramBotHandler, mqttHandler, keyboard)
	applianceInteractionHandler.KeyboardRequestHandler(telegramBotHandler)
}

func OfficeAppliancesKeyboardHandler(telegramBotHandler *TelegramBotHandler) {
	officeAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBotHandler.keyboards[OFFICE_APPLIANCES_KEYBOARD] = officeAppliancesKeyboard

	tableLampBtn := officeAppliancesKeyboard.Text(TABLE_LAMP_COMMAND)
	backBtn := officeAppliancesKeyboard.Text("⬅ Back")
	officeAppliancesKeyboard.Reply(
		officeAppliancesKeyboard.Row(tableLampBtn),
		officeAppliancesKeyboard.Row(backBtn),
	)
	usr := User{id: meId}

	telegramBotHandler.bot.Handle(OFFICE_APPLIANCES_COMMAND, func(m *tb.Message) {
		if !m.Private() {
			return
		}
		SendMessage(telegramBotHandler, usr, "Office Appliances", telegramBotHandler.keyboards[OFFICE_APPLIANCES_KEYBOARD])
	})
	telegramBotHandler.bot.Handle(&tableLampBtn, func(m *tb.Message) {
		SendMessage(telegramBotHandler, usr, "Table lamp modes", telegramBotHandler.keyboards[TABLE_LAMP_KEYBOARD])
	})
	telegramBotHandler.bot.Handle(&backBtn, func(m *tb.Message) {
		SendMessage(telegramBotHandler, usr, "/appliances", telegramBotHandler.keyboards[ALL_APPLIANCES_KEYBOARD])
	})
}

func AllAppliancesKeyboardHandler(botHandler *TelegramBotHandler) {
	allAppliancesKeyboard := &tb.ReplyMarkup{}
	botHandler.keyboards[ALL_APPLIANCES_KEYBOARD] = allAppliancesKeyboard

	officeAppliancesBtn := allAppliancesKeyboard.Text("Office appliances")
	bedRoomAppliancesBtn := allAppliancesKeyboard.Text("Bedroom appliances")

	allAppliancesKeyboard.Reply(
		allAppliancesKeyboard.Row(officeAppliancesBtn),
		allAppliancesKeyboard.Row(bedRoomAppliancesBtn),
	)

	usr := User{id: meId}

	botHandler.bot.Handle(&officeAppliancesBtn, func(m *tb.Message) {
		if !m.Private() {
			return
		}
		SendMessage(botHandler, usr, OFFICE_APPLIANCES_COMMAND, botHandler.keyboards[OFFICE_APPLIANCES_KEYBOARD])
	})
	botHandler.bot.Handle(APPLIANCES_COMMAND, func(m *tb.Message) {
		if !m.Private() {
			return
		}
		SendMessage(botHandler, usr, APPLIANCES_COMMAND, botHandler.keyboards[ALL_APPLIANCES_KEYBOARD])
	})
}
