package main

import (
	"encoding/json"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"net/http"
)

type CryptoQuery struct {
}

func (cryptoQuery *CryptoQuery) RequestData(cryptoCurrency string) *http.Response {
	request := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=" + cryptoCurrency + "&order=market_cap_desc&per_page=1&page=1&price_change_percentage=24h%2C7d%2C14d%2C30d"
	response, err := http.Get(request)
	if err != nil {
		log.Fatalln("failed to get response", err)
	}
	return response
}

func (cryptoQuery *CryptoQuery) CreateBody(response *http.Response) []byte {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

func (cryptoQuery *CryptoQuery) FilterUnwanted(body []byte) map[string]interface{} {
	var queryResult []map[string]interface{}
	err := json.Unmarshal(body, &queryResult)
	if err != nil {
		fmt.Println(err)
	}

	filteredData := make(map[string]interface{})
	filteredData["price"] = queryResult[0]["current_price"]
	filteredData["high_24h"] = queryResult[0]["high_24h"]
	filteredData["low_24h"] = queryResult[0]["low_24h"]
	filteredData["price_change_24h"] = queryResult[0]["price_change_24h"]
	filteredData["pc_change_perc_24h"] = queryResult[0]["price_change_percentage_24h"]
	filteredData["pc_change_perc_7d"] = queryResult[0]["price_change_percentage_7d_in_currency"]
	filteredData["pc_change_perc_14d"] = queryResult[0]["price_change_percentage_14d_in_currency"]
	/*filteredData["image"] = queryResult[0]["image"]*/
	return filteredData
}

func (cryptoQuery *CryptoQuery) GenerateFunctionButtons(services *ServiceContainer) map[string]*tb.Btn {
	buttons := make(map[string]*tb.Btn)

	buttons["bitcoin"] = &tb.Btn{Unique: "bitcoin", Text: "BTC"}
	buttons["ethereum"] = &tb.Btn{Unique: "ethereum", Text: "ETH"}

	for currency, btn := range buttons {

		func(btn *tb.Btn, currency string) {

			services.botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := services.botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					fmt.Println("REEEEEEROOO", err)
				}

				response := cryptoQuery.RequestData(currency)
				body := cryptoQuery.CreateBody(response)
				cryptoData := cryptoQuery.FilterUnwanted(body)
				humanReadable := CreateHumanReadable(cryptoData)

				_, err = services.botHandler.bot.Send(&me, humanReadable)
				if err != nil {
					fmt.Println(err)
				}
			})
		}(btn, currency)
	}

	return buttons
}

func (cryptoQuery *CryptoQuery) KeyboardCommands(services *ServiceContainer) {

	buttons := cryptoQuery.GenerateFunctionButtons(services)

	cryptoDataKeyboard := &tb.ReplyMarkup{}

	cryptoDataKeyboard.Inline(
		cryptoDataKeyboard.Row(*buttons["bitcoin"], *buttons["ethereum"]))

	services.botHandler.keyboards[CRYPTO_DATA_KEYBOARD] = cryptoDataKeyboard

}

func (cryptoQuery *CryptoQuery) NonKeyboardCommands(services *ServiceContainer) {
	services.botHandler.UserEvent("/crypto", "Crypto Data", CRYPTO_DATA_KEYBOARD, KBOARD)
}
