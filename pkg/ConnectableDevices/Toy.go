package connectable

type Toy struct {
	Name           string
	AvailableModes []string
	Room           string
	Id             int
	PublishTopic   string
	SubscribeTopic string
	BotCommand     string
}
