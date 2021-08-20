package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const masterChatID int64 = 558297691
var Bot, botError = tgbotapi.NewBotAPI("1763947554:AAHUWq4nR30Hj7WybEbR3ztSwm2CO2C-X4k")

func SetupBot()  {

	Bot.Debug = false

	if botError != nil {
		log.Panic(botError)
	}
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


func CreateUserReply(humanReadable string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(masterChatID, humanReadable)
}

func CreateUpdateConfig() tgbotapi.UpdateConfig{
	botUpdte := tgbotapi.NewUpdate(0)
	botUpdte.Timeout = 60
	return botUpdte
}
