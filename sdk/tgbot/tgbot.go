package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"net/http"
)

type ChatProvider interface {
	GetChat(chatId int64) (ChatInfo, error)
	SaveChatInfo(chat ChatInfo) error
	GetButton(btnId string) (Button, error)
	SaveButton(button Button) error
}

type Bot struct {
	*tgbotapi.BotAPI
	handler  func(update *Update)
	chatProv ChatProvider
}

func NewBot(token string, chatProv ChatProvider) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create TgBot")
	}
	return &Bot{BotAPI: api, chatProv: chatProv}, nil
}

func (b *Bot) StartLongPolling(handler func(update *Update)) error {
	if b.handler != nil {
		return errors.New("long polling already started")
	}
	b.handler = handler
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range b.GetUpdatesChan(u) {
		wrappedUpdate := b.WrapUpdate(update)
		b.handler(wrappedUpdate)
	}
	return nil
}

func (b *Bot) WrapUpdate(update tgbotapi.Update) *Update {
	return WrapUpdate(update, b.chatProv)
}

func (b *Bot) WrapRequest(req *http.Request) (*Update, error) {
	update, err := b.HandleUpdate(req)
	if err != nil {
		return nil, err
	}
	return WrapUpdate(*update, b.chatProv), nil
}
