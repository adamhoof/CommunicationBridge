package main

type ServiceContainer struct {
	mqtt       *MQTTHandler
	botHandler *TelegramBot
	db         *PostgreSQLHandler
}
