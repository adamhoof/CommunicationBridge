package TelegramBot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const masterChatID int64 = 558297691
var Bot, botError = tgbotapi.NewBotAPI("1763947554:AAHUWq4nR30Hj7WybEbR3ztSwm2CO2C-X4k")

func Setup()  {

	Bot.Debug = false

	if botError != nil {
		log.Panic(botError)
	}
}

func CreateHumanReadable(intrfeic interface{}) string {

	if applianceData, ok := intrfeic.([]string) ; ok{

		if applianceData != nil {

			var humanReadable string

			for i := 0; i < len(applianceData); i++ {
				humanReadable += applianceData[i] + "\n"

			}
			return humanReadable
		}
		return "Failed to set device"
	}
	return intrfeic.(string)
}

func CreateUserReply(humanReadable string) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(masterChatID, humanReadable)
}
