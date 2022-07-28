package main

import (
	connectable "RPICommandHandler/pkg/ConnectableDevices"
	"RPICommandHandler/pkg/Database"
	env "RPICommandHandler/pkg/Env"
	"RPICommandHandler/pkg/MQTT"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
)

func main() {
	env.SetEnv()

	mqttHandler := mqtts.Client{}
	postgresHandler := database.PostgresHandler{}

	var routineSyncer sync.WaitGroup

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		mqttHandler.SetupClientOptions()
		mqttHandler.CreateClient()
		mqttHandler.ConnectClient()
	}(&routineSyncer)

	routineSyncer.Add(1)
	go func(routineSyncer *sync.WaitGroup) {
		defer routineSyncer.Done()
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s",
			os.Getenv("dbHost"),
			os.Getenv("dbPort"),
			os.Getenv("dbUser"),
			os.Getenv("dbPassword"),
			os.Getenv("fads"))

		if err := postgresHandler.Connect(&psqlInfo); err != nil {
			panic(err)
		}
	}(&routineSyncer)

	routineSyncer.Wait()

	toys := make(map[string]*connectable.Toy)
	postgresHandler.PullToyData(toys)
	fmt.Println("shit")
}
