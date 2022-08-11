package main

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	database "RPICommandHandler/pkg/Database"
	env "RPICommandHandler/pkg/Env"
	response "RPICommandHandler/pkg/MQTTResponseHandlers"
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

	roomKeyboards := make(map[string]*tb.ReplyMarkup)

	rooms := postgresHandler.PullAvailableRooms()
	backBtn := tb.Btn{Text: "⬅"}
	botHandler.SendKeyboardOnButtonClick(&backBtn)

	for _, room := range rooms {
		toys := make(map[string]*connectable.Toy)
		postgresHandler.PullToyDataBasedOnRoom(toys, room)
		replyButtons := telegram.CreateReplyButtonsForToysInRoom(toys, room[0:2])
		keyboard := telegram.CreateReplyKeyboardFromButtons(replyButtons, backBtn)
		roomKeyboards[room] = keyboard

		for _, toy := range toys {
			for _, replyButton := range replyButtons {
				botHandler.SendKeyboardOnButtonClick(&replyButton, toy.Name, roomKeyboards, toy.Name)
			}
			responseHandler := response.DefaultDeviceResponseHandler(&botHandler, toy.Name) //TODO toy will have a field [handlerName string], which will then be put into a map to lookup handler with that name
			mqttClient.Subscribe(toy.SubscribeTopic, 0, responseHandler)

			inlineButtons := telegram.GenerateInlineButtonsForToy(toy)

			telegram.GenerateInlineKeyboardFromButtonsForToy(inlineButtons)
		}
	}

	telegram.CreateAllToysKeyboardUI(&botHandler, roomKeyboards)

	botHandler.StartBot()
}
