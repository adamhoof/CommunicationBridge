package main

import (
	_ "github.com/lib/pq"
	tb "gopkg.in/tucnak/telebot.v2"
	"sync"
)

func main() {

	mqttHandler := MQTTHandler {}
	postgreSQLHandler := PostgreSQLHandler{}
	telegramBotHandler := TelegramBotHandler{}

	var routineSyncer sync.WaitGroup
	mappie := make(map[string]interface{})
	mappie["Type"] = "fuck"

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		telegramBotHandler.CreateBot()
		telegramBotHandler.CreateTableLampKeyboard()
		telegramBotHandler.bot.Handle("/TableLamp", func(message *tb.Message) {
			if !message.Private() {
				return
			}
			telegramBotHandler.bot.Send(message.Sender, "Table Lamp Colors", telegramBotHandler.tableLampModeKeyboard)
		})
		telegramBotHandler.bot.Start()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		mqttHandler.SetupTLSConfig()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
		mqttHandler.SetSubscriptions()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(wg *sync.WaitGroup) {
		defer routineSyncer.Done()
		postgreSQLHandler.Connect()
		postgreSQLHandler.TestConnection()
		postgreSQLHandler.CloseConnection()
	}(&routineSyncer)

	routineSyncer.Wait()

		/*DistributeCommands(&mqttHandler.client, message)*/
}
