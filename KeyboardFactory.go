package main

import tb "gopkg.in/telebot.v3"

type KeyboardFactory struct {
}

func (factory *KeyboardFactory) createFromButtons(buttons map[string]*tb.Btn) *tb.ReplyMarkup {

	var buttonsSlice = make([]tb.Btn, len(buttons))

	i := 0
	for name, _ := range buttons {
		buttonsSlice[i] = *buttons[name]
		i++
	}

	keyboard := &tb.ReplyMarkup{ResizeKeyboard: true}
	keyboard.Inline(keyboard.Row(buttonsSlice...))

	return keyboard
}
