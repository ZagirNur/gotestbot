package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	tgbot2 "gotestbot/sdk/tgbot"
	"strconv"
)

func (v *View) Share(u *tgbot2.Update) (tgbotapi.Message, error) {
	msg := new(tgbot2.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text("Вы можете объединить свои холодильники в один общий, для этого отправьте его другому").
		AddKeyboardRow().AddButtonSwitch("Отправить холодильник другому", "share").
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) ShareInline(u *tgbot2.Update) (tgbotapi.Message, error) {
	shareBtn := v.createButton(ActionMerge, map[string]string{
		"chatId": strconv.FormatInt(u.GetChatId(), 10),
	})

	msg := tgbot2.NewInlineRequest(u.GetInlineId()).
		AddArticle(uuid.NewString(),
			"Объединить холодильники",
			"После объединения у вас будет один общий холодильник, все продукты будут добавлены в него").
		AddKeyboardRow().AddButton("Объединить наши холодильники", shareBtn.Id).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) GoToBotScreen(u *tgbot2.Update) (tgbotapi.Message, error) {

	msg := new(tgbot2.MessageBuilder).
		ChatId(u.GetChatId()).InlineId(u.GetInlineId()).
		Edit(u.IsButton()).
		Text("перейдите в чат с ботом").
		AddKeyboardRow().AddButtonUrl("К боту", "http://t.me/"+v.tg.BotSelf.UserName+"?start").
		Build()

	return logIfError(v.tg.Send(msg))
}
