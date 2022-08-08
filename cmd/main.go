package main

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	database "RPICommandHandler/pkg/Database"
	env "RPICommandHandler/pkg/Env"
	mqtt2 "RPICommandHandler/pkg/MQTT"
	telegram "RPICommandHandler/pkg/Telegram"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

func main() {
	env.SetEnv()

	mqttClient := mqtt2.MQTTClient{}
	postgresHandler := database.PostgresHandler{}
	me := telegram.User{Id: os.Getenv("telegramBotOwner")}
	botHandler := telegram.BotHandler{Owner: me}

	var routineSyncer sync.WaitGroup
	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup, client *mqtt2.MQTTClient) {
		defer routineSyncer.Done()
		options := mqtt.ClientOptions{}
		options.AddBroker(os.Getenv("mqttServer"))
		options.SetClientID(os.Getenv("mqttClientName"))
		options.SetAutoReconnect(true)
		options.SetConnectRetry(true)
		options.SetCleanSession(false)
		options.SetOrderMatters(false)
		// use options.SetTLSConfig if you are running MQTT broker on remote server instead of localhost
		client.SetOptions(options)
		client.Connect()
	}(&routineSyncer, &mqttClient)

	routineSyncer.Add(1)
	go func(syncer *sync.WaitGroup, handler database.DatabaseHandler) {
		defer syncer.Done()

		dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			os.Getenv("dbHost"),
			os.Getenv("dbPort"),
			os.Getenv("dbUser"),
			os.Getenv("dbPassword"),
			os.Getenv("dbName"))

		if err := postgresHandler.Connect(dbConnectionString); err != nil {
			panic(err)
		}
	}(&routineSyncer, &postgresHandler)
	routineSyncer.Wait()

	toys := make(map[string]*connectable.Toy)
	postgresHandler.PullToyData(toys)

	botHandler.CreateBot(os.Getenv("telegramBotToken"))
	botHandler.StartBot()
}
