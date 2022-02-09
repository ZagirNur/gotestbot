package sdk

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"strings"
)

type Command string

type Action string

type ChatStateProvider interface {
	GetChatState(int64) (ChatState, error)
	SaveChatState(ChatState) error
}

type ButtonProvider interface {
	GetButton(string) (Button, error)
	SaveButton(button Button) error
}

type ChatState struct {
	ChatId int64

	ActiveChain     string
	ActiveChainStep string
	Data            map[string]string
}

func (s *ChatState) WithChain() bool {
	return s.ActiveChain != ""
}

func (s *ChatState) StartChain(chain string) {
	s.ActiveChain = chain
}

func (s *ChatState) StartChainStep(chainStep string) {
	s.ActiveChainStep = chainStep
}

func (s *ChatState) HasChain(chain string) bool {
	return s.ActiveChain == chain
}

func (s *ChatState) GetChainStep() string {
	return s.ActiveChainStep
}

func (s *ChatState) AddData(key string, text string) {
	if s.Data == nil {
		s.Data = map[string]string{}
	}

	s.Data[key] = text
}

func (s *ChatState) GetData(key string) string {
	if s.Data == nil {
		s.Data = map[string]string{}
	}

	return s.Data[key]
}

func (s *ChatState) FinishChain() {
	s.ActiveChain = ""
	s.ActiveChainStep = ""
}

type Button struct {
	Id     string
	Action Action
	Data   map[string]string
}

func (b Button) HasAction(action Action) bool {
	return b.Action == action
}

func (b Button) GetData(key string) string {
	if b.Data == nil {
		return ""
	}

	return b.Data[key]
}

type Update struct {
	tgbotapi.Update
	stateProv ChatStateProvider
	btnProv   ButtonProvider

	chatState *ChatState
	btn       *Button
}

func NewUpdate(update tgbotapi.Update, stateProv ChatStateProvider, btnProv ButtonProvider) *Update {
	return &Update{Update: update, stateProv: stateProv, btnProv: btnProv}
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

func (u *Update) IsPlainText() bool {
	return !u.IsCommand() && u.Update.Message != nil && u.Update.Message.Text != ""
}

func (u *Update) IsButton() bool {
	return u.Update.CallbackData() != ""
}

func (u *Update) HasAction(action Action) bool {
	return u.IsButton() && u.GetButton().HasAction(action)
}

func (u *Update) HasActionOrChain(actionOrChain Action) bool {
	return u.IsButton() && u.GetButton().HasAction(actionOrChain) ||
		u.GetChatState().HasChain(string(actionOrChain))
}

func (u *Update) GetText() string {
	return u.Message.Text
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

func (u *Update) GetChatState() *ChatState {
	if u.chatState == nil {
		state, err := u.stateProv.GetChatState(u.GetChatId())
		if err != nil {
			log.Error().Err(err).Msgf("cannot find chat state")
		}
		if state.ChatId == 0 {
			state.ChatId = u.GetChatId()
		}
		u.chatState = &state
	}
	return u.chatState
}

func (u *Update) GetButton() Button {
	if u.btn == nil {
		button, err := u.btnProv.GetButton(u.CallbackData())
		if err != nil {
			log.Error().Err(err).Msgf("cannot find button %s", u.CallbackData())
		}
		u.btn = &button
	}
	return *u.btn
}

func (u *Update) FlushState() {

	go func() {

		err := u.stateProv.SaveChatState(*u.GetChatState())
		if err != nil {
			log.Error().Err(err).Msgf("cannot save chat state: %+v", u.GetChatState())
		}
	}()
}

func (u *Update) FinishChain() *Update {
	u.GetChatState().ActiveChain = ""
	u.GetChatState().ActiveChainStep = ""
	return u
}

func (u *Update) StartChain(chain string) *Update {
	u.GetChatState().ActiveChain = chain
	return u
}

func (u *Update) StartChainStep(chainStep string) *Update {
	u.GetChatState().ActiveChainStep = chainStep
	return u
}

func (u *Update) AddChainData(key string, value string) *Update {
	u.GetChatState().AddData(key, value)
	return u
}
