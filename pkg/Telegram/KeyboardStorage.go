package telegrambot

import tb "gopkg.in/telebot.v3"

type KeyboardStorage struct {
	keyboards map[string]*tb.ReplyMarkup
}

func (storage *KeyboardStorage) AddKeyboard(name string, keyboard *tb.ReplyMarkup) {
	storage.keyboards[name] = keyboard
}

func (storage *KeyboardStorage) GetKeyboardByName(name string) *tb.ReplyMarkup {
	return storage.keyboards[name]
}
