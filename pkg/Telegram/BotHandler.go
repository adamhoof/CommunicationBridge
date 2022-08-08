package telegrambot

import (
	"fmt"
	"github.com/adamhoof/GolangTypeConvertorWrapper/pkg"
	tb "gopkg.in/telebot.v3"
	"time"
)

type BotHandler struct {
	Bot   *tb.Bot
	Owner User
}

func (handler *BotHandler) CreateBot(token string) {
	var err error
	handler.Bot, err = tb.NewBot(tb.Settings{
		Token: token,
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Telegram token valid")
}

func (handler *BotHandler) StartBot() {
	handler.Bot.Start()
}

func (handler *BotHandler) OwnerVerify(id int64) bool {
	if id != typeconv.StringToInt64(handler.Owner.Id) {
		return false
	}
	return true
}

func (handler *BotHandler) SendTextMessage(owner *User, title string, message string) {

	_, err := handler.Bot.Send(owner, title, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}
