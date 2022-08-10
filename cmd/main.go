package main

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	database "RPICommandHandler/pkg/Database"
	env "RPICommandHandler/pkg/Env"
	telegram "RPICommandHandler/pkg/Telegram"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	tb "gopkg.in/telebot.v3"
	"os"
	"sync"
)

func main() {
	env.SetEnv()

	options := mqtt.ClientOptions{}
	options.AddBroker(os.Getenv("mqttServer"))
	options.SetClientID(os.Getenv("mqttClientName"))
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)
	options.SetCleanSession(false)
	options.SetOrderMatters(false)
	// use options.SetTLSConfig if you want to establish secure connection (not required on localhost, recommended when connecting to remote server)
	mqttClient := mqtt.NewClient(&options)
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

	roomKeyboards := make(map[string]*tb.ReplyMarkup) //create storage for keyboards THAT HOLD buttons for individual toys -> Office holds buttons for lamp, shades...
	/*toyKeyboards := make(map[string]*tb.ReplyMarkup)*/ // create storage for keyboards of individual toys -> Each toy has its own keyboard of commands

	toys := make(map[string]*connectable.Toy)
	postgresHandler.PullToyData(toys)

	telegram.CreateRoomOfToysKeyboard(&botHandler, roomKeyboards, toys, telegram.BedroomToysKeyboard, "b ")
	/*handler.SendKeyboardOnButtonClick(&button, toy.Name+" modes", keyboards, toy.Name)*/
	telegram.CreateRoomOfToysKeyboard(&botHandler, roomKeyboards, toys, telegram.OfficeToysKeyboard, "o ")

	telegram.CreateAllToysKeyboardUI(&botHandler, roomKeyboards)

	botHandler.StartBot()
}
