package main

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	database "RPICommandHandler/pkg/Database"
	env "RPICommandHandler/pkg/Env"
	telegram "RPICommandHandler/pkg/Telegram"
	"fmt"
	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	tb "gopkg.in/telebot.v3"
	"os"
	"sync"
)

func main() {
	env.SetEnv()

	options := pahomqtt.ClientOptions{}
	options.AddBroker(os.Getenv("mqttServer"))
	options.SetClientID(os.Getenv("mqttClientName"))
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetCleanSession(false)
	options.SetOrderMatters(false)
	// use options.SetTLSConfig if you want to establish secure connection (not required on localhost, recommended when connecting to remote server)
	mqttClient := pahomqtt.NewClient(&options)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("MQTT client connection established?: ", mqttClient.IsConnected())

	postgresHandler := database.PostgresHandler{}

	var routineSyncer sync.WaitGroup
	routineSyncer.Add(1)
	go func(syncer *sync.WaitGroup, handler database.DatabaseHandler) {
		defer syncer.Done()
		dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			os.Getenv("dbHost"),
			os.Getenv("dbPort"),
			os.Getenv("dbUser"),
			os.Getenv("dbPassword"),
			os.Getenv("dbName"))

		if err := handler.Connect(dbConnectionString); err != nil {
			panic(err)
		}
	}(&routineSyncer, &postgresHandler)
	routineSyncer.Wait()

	me := telegram.User{Id: os.Getenv("telegramBotOwner")}
	botHandler := telegram.BotHandler{Owner: me}
	botHandler.CreateBot(os.Getenv("telegramBotToken"))

	keyboards := make(map[string]*tb.ReplyMarkup)
	telegram.CreateAllToysKeyboardUI(&botHandler, keyboards)
	telegram.CreateOfficeToysKeyboardUI(&botHandler, keyboards)
	telegram.CreateBedroomToysKeyboardUI(&botHandler, keyboards)

	toys := make(map[string]*connectable.Toy)
	postgresHandler.PullToyData(toys)

	botHandler.StartBot()
}
