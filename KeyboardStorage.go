package main

import tb "gopkg.in/telebot.v3"

type KeyboardStorage struct {
	keyboards map[string]*tb.ReplyMarkup
}

func (storage *KeyboardStorage) unlock() {
	storage.keyboards = make(map[string]*tb.ReplyMarkup)
}

func (storage *KeyboardStorage) store(keyboardName string, keyboard *tb.ReplyMarkup) {
	storage.keyboards[keyboardName] = keyboard
}

func (storage *KeyboardStorage) getKeyboardByName(name string) *tb.ReplyMarkup {
	return storage.keyboards[name]
}
