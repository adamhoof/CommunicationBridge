package telegram

import tb "gopkg.in/telebot.v3"

type KeyboardStorage struct {
	Keyboards map[string]*tb.ReplyMarkup
}

func (storage *KeyboardStorage) AddKeyboard(name string, keyboard *tb.ReplyMarkup) {
	storage.Keyboards[name] = keyboard
}

func (storage *KeyboardStorage) GetKeyboardByName(name string) *tb.ReplyMarkup {
	return storage.Keyboards[name]
}
