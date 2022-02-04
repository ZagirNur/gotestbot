package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gotestbot/internal/view"
	"gotestbot/sdk"
	"strconv"
)

type BotApp struct {
	view  *view.View
	repos *Repository
}

func NewBotApp(view *view.View, repos *Repository) *BotApp {
	return &BotApp{view: view, repos: repos}
}

func (b *BotApp) Handle(update tgbotapi.Update) error {
	u := sdk.NewUpdate(update, b.repos)

	if u.HasCommand("/start") {
		go b.view.StartView(u)
		return nil
	}

	if u.HasButtonAction("REGISTRATION") || u.ChatState().HasChain("REGISTRATION") {
		b.HandleRegistration(u)
	}

	if u.HasButtonAction("LEAVE") {
		b.repos.DeleteUser(u.GetChatId())
		go b.view.StartView(u)
	}

	return nil
}

func (b BotApp) HandleRegistration(u *sdk.Update) {
	chatState := u.ChatState()
	switch chatState.GetChainStep() {
	case "":
		chatState.StartChain("REGISTRATION")
		chatState.StartChainStep("NAME")
		u.FlushState()
		go b.view.RegistrationEnterNameView(u)
	case "NAME":
		if u.IsPlainText() {
			chatState.AddData("name", u.GetText())
			chatState.StartChainStep("AGE")
			u.FlushState()
			go b.view.RegistrationEnterAgeView(u, u.GetText())
		}
	case "AGE":
		if u.IsPlainText() {
			name := chatState.GetData("name")
			age, err := strconv.Atoi(u.GetText())
			if err != nil {
				b.view.RegistrationEnterAgeErrorView(u, name)
				return
			}
			b.repos.SaveUser(sdk.User{Id: u.GetChatId(), Name: name, Age: age})
			chatState.FinishChain()
			u.FlushState()
			go b.view.StartView(u)
		}
	}

}
