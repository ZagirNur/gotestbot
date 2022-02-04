package sdk

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type ChatStateProvider interface {
	GetChatState(int64) ChatState
	SaveChatState(ChatState)
	GetButton(string) Button
}

type ChatState struct {
	ChatId int64

	activeChain     string
	activeChainStep string
	data            map[string]string
}

func (s *ChatState) WithChain() bool {
	return s.activeChain != ""
}

func (s *ChatState) StartChain(chain string) {
	s.activeChain = chain
}

func (s *ChatState) StartChainStep(chainStep string) {
	s.activeChainStep = chainStep
}

func (s *ChatState) HasChain(chain string) bool {
	return s.activeChain == chain
}

func (s *ChatState) GetChainStep() string {
	return s.activeChainStep
}

func (s *ChatState) AddData(key string, text string) {
	if s.data == nil {
		s.data = map[string]string{}
	}

	s.data[key] = text
}

func (s *ChatState) GetData(key string) string {
	if s.data == nil {
		s.data = map[string]string{}
	}

	return s.data[key]
}

func (s *ChatState) FinishChain() {
	s.activeChain = ""
	s.activeChainStep = ""
}

type Button struct {
	Id     string
	Action string
}

type User struct {
	Id   int64
	Name string
	Age  int
}

func (b Button) HasAction(action string) bool {
	return b.Action == action
}

type Update struct {
	tgbotapi.Update
	chatState     *ChatState
	btn           *Button
	stateProvider ChatStateProvider
}

func NewUpdate(update tgbotapi.Update, stateProvider ChatStateProvider) *Update {
	return &Update{Update: update, stateProvider: stateProvider}
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

func (u *Update) HasButtonAction(action string) bool {
	return u.IsButton() && u.stateProvider.GetButton(u.CallbackData()).HasAction(action)
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
	return u.Message.MessageID
}

func (u *Update) ChatState() *ChatState {
	if u.chatState == nil {
		state := u.stateProvider.GetChatState(u.GetChatId())
		u.chatState = &state
	}
	return u.chatState
}

func (u *Update) GetButton() Button {
	if u.btn == nil {
		btn := u.stateProvider.GetButton(u.CallbackData())
		u.btn = &btn
	}
	if u.btn == nil {
		return Button{}
	}
	return *u.btn
}

func (u *Update) FlushState() {
	u.stateProvider.SaveChatState(*u.ChatState())
}

func (u *Update) FinishChain() {

}

type Message struct {
	tgbotapi.Message
}
