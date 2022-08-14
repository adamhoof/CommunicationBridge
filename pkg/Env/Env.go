package env

import (
	"fmt"
	"os"
)

func SetEnv() {
	err := os.Setenv("dbHost", "localhost")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("dbPort", "5432")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("dbUser", "adamhoof")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("dbPassword", "medprodsdb")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("dbName", "medunkaproducts")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("mqttServer", "tcp://raspberrypi.local:1883")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("mqttClientName", "RPIBridge")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("telegramBotOwner", "558297691")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
	err = os.Setenv("telegramBotToken", "1914152683:AAEaxUTD05d4pMluKCeDWeBwPdQQC0KOGwE")
	if err != nil {
		fmt.Println("unable to set env for: ", err)
	}
}
