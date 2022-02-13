package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tgbot2 "gotestbot/sdk/tgbot"
)

func (v *View) AddProductName(u *tgbot2.Update) (tgbotapi.Message, error) {

	builder := new(tgbot2.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text("Введите название продукта")

	return logIfError(v.tg.Send(builder.Build()))
}

func (v *View) AddProductDate(prefix string, u *tgbot2.Update) (tgbotapi.Message, error) {
	builder := new(tgbot2.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text(prefix + "Введите срок годности в виде дд.мм.гггг, например 19.02.2022")

	return logIfError(v.tg.Send(builder.Build()))
}
