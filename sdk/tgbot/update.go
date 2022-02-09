package tgbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

type Update struct {
	tgbotapi.Update
	chatProv ChatProvider

	chat *ChatInfo
	btn  *Button
}

func WrapUpdate(update tgbotapi.Update, stateProv ChatProvider) *Update {
	return &Update{Update: update, chatProv: stateProv}
}

func (u *Update) GetChatId() int64 {
	if u.Message != nil && u.Message.Chat != nil {
		return u.Message.Chat.ID
	}
	if u.CallbackQuery != nil {
		return u.CallbackQuery.From.ID
	}
	return 0
}

func (u *Update) GetMessageId() int {
	if u.IsButton() {
		return u.CallbackQuery.Message.MessageID
	} else if u.Message != nil {
		return u.Message.MessageID
	}
	return 0
}

func (u *Update) HasText(text string) bool {
	return text == u.Update.Message.Text
}

func (u *Update) HasCommand(text string) bool {
	return u.IsCommand() && text == u.Update.Message.Text
}

func (u *Update) IsCommand() bool {
	return u.Update.Message != nil &&
		strings.Contains(u.Update.Message.Text, "/")
}

//Button

func (u *Update) IsPlainText() bool {
	return !u.IsCommand() && u.Update.Message != nil && u.Update.Message.Text != ""
}

func (u *Update) GetText() string {
	return u.Message.Text
}

func (u *Update) IsButton() bool {
	return u.Update.CallbackData() != ""
}

func (u *Update) GetButton() Button {
	if u.btn == nil {
		button, err := u.chatProv.GetButton(u.CallbackData())
		if err != nil {
			log.Error().Err(err).Msgf("cannot find button %s", u.CallbackData())
		}
		u.btn = &button
	}
	return *u.btn
}

func (u *Update) HasAction(action Action) bool {
	return u.IsButton() && u.GetButton().HasAction(action)
}

func (u *Update) HasActionOrChain(actionOrChain Action) bool {
	return u.IsButton() && u.GetButton().HasAction(actionOrChain) ||
		u.GetChatInfo().ActiveChain == string(actionOrChain)
}

// ChatInfo

func (u *Update) GetChatInfo() *ChatInfo {
	if u.chat == nil {
		chat, err := u.chatProv.GetChat(u.GetChatId())
		if err != nil {
			log.Error().Err(err).Msgf("cannot find chat chat")
		}
		u.chat = &chat
	}

	if u.chat.ChatId == 0 {
		u.chat.ChatId = u.GetChatId()
	}
	if u.chat.ChainData == nil {
		u.chat.ChainData = Data{}
	}

	return u.chat
}

func (u *Update) FlushChatInfo() {
	err := u.chatProv.SaveChatInfo(*u.GetChatInfo())
	if err != nil {
		log.Error().Err(err).Msgf("cannot save chat info: %+v", u.GetChatInfo())
	}
}

func (u *Update) StartChain(chain string) *Update {
	u.GetChatInfo().ActiveChain = chain
	return u
}

func (u *Update) StartChainStep(chainStep string) *Update {
	u.GetChatInfo().ActiveChainStep = chainStep
	return u
}

func (u *Update) GetChainStep() string {
	return u.GetChatInfo().ActiveChainStep
}

func (u *Update) AddChainData(key string, value string) *Update {
	u.GetChatInfo().ChainData[key] = value
	return u
}

func (u *Update) GetChainData(key string) string {
	return u.GetChatInfo().ChainData[key]
}

func (u *Update) FinishChain() *Update {
	u.GetChatInfo().ActiveChain = ""
	u.GetChatInfo().ActiveChainStep = ""
	u.GetChatInfo().ChainData = map[string]string{}
	return u
}
