package connectable

type Toy struct {
	Name           string   `json:"name"`
	AvailableModes []string `json:"availableModes"`
	PublishTopic   string   `json:"subscribeTopic"`
	SubscribeTopic string   `json:"publishTopic"`
	BotCommand     string
}
