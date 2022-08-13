package telegram

import (
	"fmt"
	"github.com/adamhoof/GolangTypeConvertorWrapper/pkg"
	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func (handler *BotHandler) HandleCommand(command string, h tb.HandlerFunc) {
	handler.Bot.Handle(command, h)
}

func (handler *BotHandler) HandleButtonClick(btn *tb.Btn, h tb.HandlerFunc) {
	handler.Bot.Handle(btn, h)
}

func (handler *BotHandler) SendCommandViaMQTT(command string, publishTopic string, client mqtt.Client) tb.HandlerFunc {
	return func(c tb.Context) error {
		err := handler.Bot.Respond(c.Callback(), &tb.CallbackResponse{})
		if err != nil {
			return err
		}
		client.Publish(publishTopic, 0, true, command)
		return nil
	}
}

func (handler *BotHandler) SendKeyboard(title string, keyboards map[string]*tb.ReplyMarkup, keyboardName string) tb.HandlerFunc {
	return func(c tb.Context) (err error) {
		if !handler.OwnerVerify(c.Message().Sender.ID) {
			handler.SendTextMessage(&handler.Owner, "failed to verify owner")
			fmt.Println("owner not verified", err)
		}
		_, err = handler.Bot.Send(&handler.Owner, title, keyboards[keyboardName])
		if err != nil {
			fmt.Println("could not send keyboard on command click:", err)
		}
		return nil
	}
}

func (handler *BotHandler) SendCommandsList(listOfCommands []tb.Command) tb.HandlerFunc {
	return func(c tb.Context) error {
		if !handler.OwnerVerify(c.Message().Sender.ID) {
			handler.SendTextMessage(&handler.Owner, "failed to verify owner")
			fmt.Println("owner not verified")
		}
		var commandsListAsString string = "Commands list: \n\n"
		for _, command := range listOfCommands {
			commandsListAsString += fmt.Sprintf("%s\n\n", command.Text)
		}
		handler.SendTextMessage(&handler.Owner, commandsListAsString)
		return nil
	}
}
