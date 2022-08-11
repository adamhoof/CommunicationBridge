package connectable

type Toy struct {
	Name           string
	AvailableModes []string
	Room           string
	Id             int
	PublishTopic   string
	SubscribeTopic string
}

/*
func AwakenButtons(toy *Toy, buttons map[string]*tb.Btn, client mqtt.Client) {

	for mode, btn := range buttons {

		func(btn *tb.Btn, mode string) {

			services.botHandler.bot.Handle(btn, func(c tb.Context) error {
				err := services.botHandler.bot.Respond(c.Callback(), &tb.CallbackResponse{})
				if err != nil {
					return err
				}
				client.Publish(toy.PublishTopic, 0, true, mode)
				return nil
			})
		}(btn, mode)
	}
}
*/
