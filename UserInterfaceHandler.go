package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	OFFICE_APPLIANCES_COMMAND = "/officeappliances"
	OFFICE_APPLIANCES_KEYBOARD = "officeAppliances"
	TABLE_LAMP_COMMAND = "/tablelamp"
	TABLE_LAMP_KEYBOARD = "tableLamp"
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

func RoomAppliancesKeyboardRequestHandler(telegramBotHandler *TelegramBotHandler) {
	roomAppliancesKeyboard := &tb.ReplyMarkup{}
	telegramBotHandler.keyboards[OFFICE_APPLIANCES_KEYBOARD] = roomAppliancesKeyboard

	tableLampBtn := roomAppliancesKeyboard.Text(TABLE_LAMP_COMMAND)
	roomAppliancesKeyboard.Reply(
		roomAppliancesKeyboard.Row(tableLampBtn),
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
}
