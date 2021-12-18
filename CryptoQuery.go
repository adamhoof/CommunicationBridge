package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tb "gopkg.in/tucnak/telebot.v2"
)

type CryptoQuery struct {
}

const (
	CRYPTO_QUERY_KBOARD = "cryptoQuery"
	cryptoQueryPub      = "cryptoquery/rpiCommand"
	cryptoQuerySub      = "cryptoquery/toyReply"
)

func (cryptoQuery *CryptoQuery) Name() string {
	return "cryptoquery"
}

func (cryptoQuery *CryptoQuery) MQTTCommandHandler(services *ServiceContainer) (handler mqtt.MessageHandler, topic string) {

	handler = func(client mqtt.Client, message mqtt.Message) {

		func() {
			queryResultMap := make(map[string]interface{})
			err := json.Unmarshal(message.Payload(), &queryResultMap)
			if err != nil {
				fmt.Println("unable to unmarshal crypto data", err)
			}
			humanReadable := CreateHumanReadable(queryResultMap)
			_, err = services.botHandler.bot.Send(&me, humanReadable)
			if err != nil {
				fmt.Println("unable to send crypto query data", err)
			}
		}()
	}
	return handler, cryptoQuerySub
}

func (cryptoQuery *CryptoQuery) GenerateFunctionButtons(services *ServiceContainer) map[string]*tb.Btn {
	buttons := make(map[string]*tb.Btn)

	buttons["bitcoin"] = &tb.Btn{Unique: "bitcoin", Text: "BTC"}
	buttons["ethereum"] = &tb.Btn{Unique: "ethereum", Text: "ETH"}
	buttons["dogecoin"] = &tb.Btn{Unique: "dogecoin", Text: "DOGE"}

	for currency, btn := range buttons {

		func(btn *tb.Btn, currency string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					fmt.Println("Invalid crypto button", err)
				}
				services.mqtt.PublishText(cryptoQueryPub, currency)
			})
		}(btn, currency)
	}
	return buttons
}

func (cryptoQuery *CryptoQuery) Kboard(services *ServiceContainer) {

	buttons := cryptoQuery.GenerateFunctionButtons(services)

	cryptoDataKeyboard := &tb.ReplyMarkup{}

	cryptoDataKeyboard.Inline(
		cryptoDataKeyboard.Row(*buttons["bitcoin"], *buttons["ethereum"], *buttons["dogecoin"]))

	services.botHandler.keyboards[CRYPTO_QUERY_KBOARD] = cryptoDataKeyboard

}

func (cryptoQuery *CryptoQuery) TextCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/crypto", "Crypto Data", CRYPTO_QUERY_KBOARD, KBOARD)
}
