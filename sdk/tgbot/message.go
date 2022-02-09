package tgbot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type MessageBuilder struct {
	editMessage bool
	chatId      int64
	messageId   int
	text        string
	keyboard    [][]tgbotapi.InlineKeyboardButton
}

func (b *MessageBuilder) EditMessageTextAndMarkup(chatId int64, messageId int) *MessageBuilder {
	b.chatId = chatId
	b.messageId = messageId
	b.editMessage = true
	return b
}

func (b *MessageBuilder) NewMessage(chatId int64) *MessageBuilder {
	b.chatId = chatId
	b.editMessage = false
	return b
}

func (b *MessageBuilder) Message(chatId int64, messageId int) *MessageBuilder {
	if messageId == 0 {
		return b.NewMessage(chatId)
	} else {
		return b.EditMessageTextAndMarkup(chatId, messageId)
	}
}

func (b *MessageBuilder) Text(text string) *MessageBuilder {
	b.text = text
	return b
}

func (b *MessageBuilder) Edit(editMessage bool) *MessageBuilder {
	b.editMessage = editMessage
	return b
}

func (b *MessageBuilder) AddKeyboardRow() *MessageBuilder {
	b.keyboard = append(b.keyboard, []tgbotapi.InlineKeyboardButton{})
	return b
}

func (b *MessageBuilder) AddButton(text, callbackData string) *MessageBuilder {
	b.keyboard[len(b.keyboard)-1] = append(b.keyboard[len(b.keyboard)-1],
		tgbotapi.InlineKeyboardButton{Text: text, CallbackData: &callbackData})
	return b
}

func (b *MessageBuilder) Build() tgbotapi.Chattable {
	if b.editMessage {
		kb := b.getKeyboard()
		var msg tgbotapi.Chattable
		if len(kb) > 0 {
			m := tgbotapi.NewEditMessageTextAndMarkup(
				b.chatId, b.messageId, b.text, tgbotapi.NewInlineKeyboardMarkup(kb...))
			m.ParseMode = tgbotapi.ModeHTML
			msg = m
		} else {
			m := tgbotapi.NewEditMessageText(b.chatId, b.messageId, b.text)
			m.ParseMode = tgbotapi.ModeHTML
			msg = m
		}

		return msg
	} else {
		msg := tgbotapi.NewMessage(b.chatId, b.text)
		keyboard := b.getKeyboard()
		if len(keyboard) > 0 {
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
		}
		msg.ParseMode = tgbotapi.ModeHTML
		return msg
	}
}

func (b *MessageBuilder) getKeyboard() [][]tgbotapi.InlineKeyboardButton {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, buttons := range b.keyboard {
		if len(buttons) > 0 {
			keyboard = append(keyboard, buttons)
		}
	}
	return keyboard
}
