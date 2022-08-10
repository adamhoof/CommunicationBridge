package telegram

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

func (handler *BotHandler) SendTextMessage(owner *User, message string) {

	_, err := handler.Bot.Send(owner, message)
	if err != nil {
		fmt.Println("Failed to send message", err)
	}
}

func (handler *BotHandler) SendKeyboardOnButtonClick(button *tb.Btn, title string, keyboards map[string]*tb.ReplyMarkup, keyboardName string) {
	handler.Bot.Handle(button, func(c tb.Context) (err error) {
		if !handler.OwnerVerify(c.Message().Sender.ID) {
			handler.SendTextMessage(&handler.Owner, "failed to verify owner")
			fmt.Println("owner not verified", err)

		}
		_, err = handler.Bot.Send(&handler.Owner, title, keyboards[keyboardName])
		if err != nil {
			fmt.Println("could not send keyboard on button click:", err)
		}
		return nil
	})
}
