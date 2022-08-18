package connectable

type Toy struct {
	Name           string   `json:"name"`
	IpAddress      string   `json:"ip"`
	AvailableModes []string `json:"availableModes"`
	Id             int
	PublishTopic   string `json:"subscribeTopic"`
	SubscribeTopic string `json:"publishTopic"`
	BotCommand     string
}
