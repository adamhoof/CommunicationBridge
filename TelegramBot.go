package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"sync"
	"time"
)

type TelegramBotHandler struct {
	bot *tb.Bot
	tableLampModeKeyboard *tb.ReplyMarkup
}

type User struct {
	userId string
}

var me = User{userId: "558297691"}

func (user *User) Recipient() string {
	return user.userId
}

func (botHandler *TelegramBotHandler) CreateBot() {
	var err error
	botHandler.bot, err = tb.NewBot(tb.Settings{
		Token:  "1914152683:AAF4r5URK9fCoJicXsCADukXuiTQSYM--U8",
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		panic(err)
	}
}

func (botHandler *TelegramBotHandler) GenerateButtons() map[string]*tb.Btn{
	
		tableLampModes := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		m := make(map[string]*tb.Btn)

		m["white"] = &tb.Btn{ Unique: "white", Text: "⬜", Data: "tlw"}
		m["yellow"] = &tb.Btn{ Unique: "yellow", Text: "\U0001F7E8", Data: "tly"}
		m["red"] =  &tb.Btn{ Unique: "red", Text: "\U0001F7E5", Data: "tlr"}
		m["off"] = &tb.Btn{ Unique: "off", Text: "🚫", Data: "tlo"}
	
	tableLampModes.Inline(
		tableLampModes.Row(*m["white"], *m["yellow"], *m["red"], *m["off"]),
	)
	botHandler.tableLampModeKeyboard = tableLampModes
	return m
}

func (botHandler *TelegramBotHandler) TableLampKeyboardRequestHandler(){
	botHandler.bot.Handle("/tablelamp", func(message *tb.Message) {
		if !message.Private() {
			return
		}
		_, err := botHandler.bot.Send(message.Sender, "Table Lamp Modes", botHandler.tableLampModeKeyboard)
		if err != nil {
			panic(err)
		}
	})
}

func (botHandler *TelegramBotHandler) TableLampActionHandlers(mqttHandler *MQTTHandler, buttons map[string]*tb.Btn){
	botHandler.TableLampKeyboardRequestHandler()

	var routineSyncer sync.WaitGroup

	for color, btn := range buttons {

		routineSyncer.Add(1)
		go func(color string, btn *tb.Btn, routineSyncer *sync.WaitGroup) {
				defer routineSyncer.Done()
				botHandler.bot.Handle(btn, func(c *tb.Callback) {
				err := botHandler.bot.Respond(c, &tb.CallbackResponse{})
				if err != nil {
					return
				}

				var payload string

				switch color {
				case "white": payload = TableLampWhiteUpdate
				case "yellow": payload = TableLampYellowUpdate
				case "red": payload = TableLampRedUpdate
				case "off": payload = TableLampOffUpdate
				default:
					panic("Unknown button")
					return
				}
				mqttHandler.PublishUpdate(TableLampPub, payload)
			})
		}(color, btn, &routineSyncer)
	}
	routineSyncer.Wait()
}

func CreateHumanReadable(applianceDataMap map[string]interface{}) string {

	var humanReadable string

		if applianceDataMap != nil {

			for key, value := range applianceDataMap {
				humanReadable += key + ": " + value.(string) + "\n"
			}
			return humanReadable
		}
		return "map iterating yeeted"
	}

func (botHandler *TelegramBotHandler) StartBot()  {
	botHandler.bot.Start()
}

func SendMessage(bot *tb.Bot, message string) {
	_, err := bot.Send(&me, message)
		if err != nil {
		panic(err)
	}
}
