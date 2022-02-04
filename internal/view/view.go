package view

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gotestbot/sdk"
)

type ButtonProvider interface {
	SaveButton(button sdk.Button) sdk.Button
}

type UserProvider interface {
	GetUser(int64) sdk.User
}

type View struct {
	btnProv  ButtonProvider
	userProv UserProvider

	tg *tgbotapi.BotAPI
}

func NewView(btnProv ButtonProvider, userProv UserProvider, tg *tgbotapi.BotAPI) *View {
	return &View{btnProv: btnProv, userProv: userProv, tg: tg}
}

func (v *View) StartView(u *sdk.Update) (tgbotapi.Message, error) {
	user := v.userProv.GetUser(u.GetChatId())
	if user == (sdk.User{}) {
		tbMsg := tgbotapi.NewMessage(u.GetChatId(), "Добро пожаловать!\n\nНе узнаю вас, пожалуйста зарегистрируйтесь=)")
		tbMsg.ParseMode = tgbotapi.ModeMarkdown
		tbMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Зарегистрироваться", CallbackData: &v.CreateButton("REGISTRATION").Id},
			},
		)
		return v.tg.Send(tbMsg)
	} else {
		msg := fmt.Sprintf("Привет %s! Я тебя сразу узнал) \n\n Тебе %s", user.Name, getAgeText(user.Age))
		tbMsg := tgbotapi.NewMessage(u.GetChatId(), msg)
		tbMsg.ParseMode = tgbotapi.ModeMarkdown
		tbMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Удалиться", CallbackData: &v.CreateButton("LEAVE").Id},
			},
		)
		return v.tg.Send(tbMsg)
	}

}

func getAgeText(age int) interface{} {
	return fmt.Sprintf("%d лет", age) //todo сделать красиво
}

func (v *View) RegistrationEnterNameView(u *sdk.Update) {
	tbMsg := tgbotapi.NewMessage(u.GetChatId(), "Регистрируемся. Первым делом напиши мне твое имя.")
	_, _ = v.tg.Send(tbMsg)
}

func (v *View) RegistrationEnterAgeView(u *sdk.Update, name string) {
	tbMsg := tgbotapi.NewMessage(u.GetChatId(), fmt.Sprintf("Очень приятно, %s. Теперь напиши мне свой возраст.", name))
	_, _ = v.tg.Send(tbMsg)
}

func (v *View) RegistrationEnterAgeErrorView(u *sdk.Update, name string) {
	tbMsg := tgbotapi.NewMessage(u.GetChatId(), fmt.Sprintf("%s, так не пойдет, введи возраст одним числом", name))
	_, _ = v.tg.Send(tbMsg)
}
