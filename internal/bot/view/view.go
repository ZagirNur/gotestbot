package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/sdk"
	"gotestbot/internal/service/model"
)

type UserProvider interface {
	GetUser(int64) (model.User, error)
}

type ProductProvider interface {
	GetProductsByUserId(userId int64) ([]model.Product, error)
}

type View struct {
	btnProv  sdk.ButtonProvider
	userProv UserProvider
	prodProv ProductProvider

	tg *tgbotapi.BotAPI
}

func NewView(btnProv sdk.ButtonProvider, userProv UserProvider, prodProv ProductProvider, tg *tgbotapi.BotAPI) *View {
	return &View{btnProv: btnProv, userProv: userProv, prodProv: prodProv, tg: tg}
}

func (v *View) StartView(u *sdk.Update) (tgbotapi.Message, error) {

	addBtn := v.CreateButton(ActionAddProduct, nil)
	showBtn := v.CreateButton(ActionShowProducts, nil)

	msg := new(sdk.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text("Добро пожаловать!\n\nЭто холодильник.").
		AddKeyboardRow().AddButton("Добавить", addBtn.Id).AddButton("Просмотр", showBtn.Id).
		Build()

	return logIfError(v.tg.Send(msg))
}

func (v *View) ShowProductView(prefix string, u *sdk.Update) (tgbotapi.Message, error) {
	products, err := v.prodProv.GetProductsByUserId(u.GetChatId())
	if err != nil {
		log.Error().Err(err).Msgf("unable to get products for userId: %d", u.GetChatId())
	}

	builder := new(sdk.MessageBuilder).
		Message(u.GetChatId(), u.GetMessageId()).
		Edit(u.IsButton()).
		Text(prefix + "Продукты в холодильнике")

	for _, product := range products {
		prodBtn := v.CreateButton("PRODUCT", map[string]string{"productId": product.Id})
		delBtn := v.CreateButton(ActionDeleteProduct, map[string]string{"productId": product.Id})

		builder.AddKeyboardRow().
			AddButton(product.Name, prodBtn.Id).
			AddButton(product.ExpirationDate.Format("02.01.2006"), prodBtn.Id).
			AddButton("Удалить", delBtn.Id)
	}

	startBtn := v.CreateButton(ActionStart, nil)
	addBtn := v.CreateButton(ActionAddProduct, nil)
	builder.AddKeyboardRow().
		AddButton("Назад", startBtn.Id).
		AddButton("Добавить", addBtn.Id)

	return logIfError(v.tg.Send(builder.Build()))
}

func (v *View) AddProductName(u *sdk.Update) (tgbotapi.Message, error) {

	builder := new(sdk.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text("Введите название продукта")

	return logIfError(v.tg.Send(builder.Build()))
}

func (v *View) AddProductDate(prefix string, u *sdk.Update) (tgbotapi.Message, error) {
	builder := new(sdk.MessageBuilder).
		NewMessage(u.GetChatId()).
		Text(prefix + "Введите срок годности в виде дд.мм.гггг, например 19.02.2022")

	return logIfError(v.tg.Send(builder.Build()))
}
