package main

import (
	connectable "CommunicationBridge/pkg/ConnectableDevices"
	deviceresponse "CommunicationBridge/pkg/ConnectableDevicesResponseHandlers"
	database "CommunicationBridge/pkg/Database"
	env "CommunicationBridge/pkg/Env"
	telegram "CommunicationBridge/pkg/Telegram"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	tb "gopkg.in/telebot.v3"
	"os"
	"sync"
)

func main() {
	env.SetEnv()

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

	toys := make(map[string]*connectable.Toy)
	postgresHandler.PullToyData(toys)
	keyboards := make(map[string]*tb.ReplyMarkup)

	mqttClient.Subscribe("/boot", 0, deviceresponse.OnNewDeviceBoot(&postgresHandler, &botHandler, mqttClient, &keyboards))

	for _, toy := range toys {
		keyboard := telegram.GenerateKeyboardWithButtonsHandlersForToy(&botHandler, mqttClient, toy)
		keyboards[toy.Name] = keyboard
		botHandler.HandleCommand(toy.BotCommand, botHandler.SendKeyboard(toy.Name, keyboards, toy.Name))
		mqttClient.Subscribe(toy.SubscribeTopic, 0, deviceresponse.Default(&botHandler, toy.Name))
	}

	botHandler.HandleCommand(telegram.AllRoomsCommand, botHandler.SendCommandsList(telegram.RoomCommands()))
	botHandler.HandleCommand(telegram.OfficeToysCommand, botHandler.SendCommandsList(telegram.OfficeToysCommands()))
	botHandler.HandleCommand(telegram.BedroomToysCommand, botHandler.SendCommandsList(telegram.BedroomToysCommands()))

	routineSyncer.Add(1)
	go func() { botHandler.StartBot() }()

	err := botHandler.Bot.SetCommands(telegram.RoomCommands())
	if err != nil {
		fmt.Println(err)
	}
	routineSyncer.Wait()
}
