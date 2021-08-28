package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type TelegramBotHandler struct {
	bot *tb.Bot
	tableLampModeKeyboard *tb.ReplyMarkup
}

const masterChatID int64 = 558297691

func (telegramBotHandler* TelegramBotHandler) CreateBot() {
	var err error
	telegramBotHandler.bot, err = tb.NewBot(tb.Settings{
		Token:  "1914152683:AAF4r5URK9fCoJicXsCADukXuiTQSYM--U8",
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}
}

func (telegramBotHandler* TelegramBotHandler) CreateTableLampKeyboard(){

	var (
		tableLampModes = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		btnWhite  = tableLampModes.Data("⬜", "white", "white")
		btnYellow = tableLampModes.Data("\U0001F7E8", "yellow", "yellow")
		btnRed = tableLampModes.Data("\U0001F7E5", "red", "red")
		btnOff = tableLampModes.Data("⬛", "off", "off")

	)
	tableLampModes.Inline(
		tableLampModes.Row(btnWhite, btnYellow, btnRed, btnOff),
	)
	telegramBotHandler.tableLampModeKeyboard = tableLampModes
}

func CreateHumanReadable(applianceDataMap map[string]interface{}) string {

	var humanReadable string

		if applianceDataMap != nil {

			for key, value := range applianceDataMap {
				humanReadable += key + ": " + value.(string) + "\n"
			}
			return humanReadable
		}
		return "map iterating yeeted"
	}
