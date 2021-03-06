package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/service/model"
	tgbot2 "gotestbot/sdk/tgbot"
)

type ProductProvider interface {
	GetProductsByChatId(chatId int64) ([]model.Product, error)
}

type View struct {
	chatProv tgbot2.ChatProvider
	prodProv ProductProvider

	tg *tgbot2.Bot
}

func NewView(btnProv tgbot2.ChatProvider, prodProv ProductProvider, tg *tgbot2.Bot) *View {
	return &View{chatProv: btnProv, prodProv: prodProv, tg: tg}
}

func (v *View) StartView(u *tgbot2.Update) (tgbotapi.Message, error) {

	addBtn := v.createButton(ActionAddProduct, nil)
	showBtn := v.createButton(ActionShowProducts, nil)

	msg := new(tgbot2.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text("Добро пожаловать!\n\nЭто холодильник.").
		AddKeyboardRow().AddButton("Добавить", addBtn.Id).AddButton("Просмотр", showBtn.Id).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) ShowProductView(prefix string, u *tgbot2.Update) (tgbotapi.Message, error) {
	products, err := v.prodProv.GetProductsByChatId(u.GetChatId())
	if err != nil {
		log.Error().Err(err).Msgf("unable to get products for userId: %d", u.GetChatId())
	}

	builder := new(tgbot2.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text(prefix + "Продукты в холодильнике")

	for _, product := range products {
		prodBtn := v.createButton("PRODUCT", map[string]string{"productId": product.Id.String()})
		delBtn := v.createButton(ActionDeleteProduct, map[string]string{"productId": product.Id.String()})

		builder.AddKeyboardRow().
			AddButton(product.Name, prodBtn.Id).
			AddButton(product.ExpirationDate.Format("02.01.2006"), prodBtn.Id).
			AddButton("Удалить", delBtn.Id)
	}

	startBtn := v.createButton(ActionStart, nil)
	addBtn := v.createButton(ActionAddProduct, nil)
	builder.AddKeyboardRow().
		AddButton("Назад", startBtn.Id).
		AddButton("Добавить", addBtn.Id)

	return logIfError(v.tg.Send(builder.Build()))
}

func (v *View) ErrorMessage(u *tgbot2.Update, text string) (tgbotapi.Message, error) {
	c := &tgbotapi.CallbackConfig{
		CallbackQueryID: u.CallbackQuery.ID,
		Text:            text,
		ShowAlert:       true,
	}
	return logIfError(v.tg.Send(c))
}
