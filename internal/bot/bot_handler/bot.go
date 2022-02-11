package bot_handler

import (
	"github.com/google/uuid"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/service/model"
	"gotestbot/sdk/tgbot"
	"time"
)

type ProductProvider interface {
	GetProductsByChatId(chatId int64) ([]model.Product, error)
	SaveProduct(product model.Product) error
	DeleteProduct(productId string) error
}

type BotApp struct {
	view     *view.View
	prodProv ProductProvider
}

func NewBotApp(view *view.View, prodProv ProductProvider) *BotApp {
	return &BotApp{view: view, prodProv: prodProv}
}

func (b *BotApp) Handle(u *tgbot.Update) {
	switch {
	case u.HasCommand("/start") || u.HasAction(view.ActionStart):
		u.FinishChain().FlushChatInfo()
		_, _ = b.view.StartView(u)

	case u.HasActionOrChain(view.ActionAddProduct):
		b.HandleAddProduct(u)

	case u.HasAction(view.ActionDeleteProduct):
		err := b.prodProv.DeleteProduct(u.GetButton().GetData("productId"))
		if err != nil {
			_, _ = b.view.ShowProductView("Ошибка удаления продукта.\n", u)
			return
		}
		_, _ = b.view.ShowProductView("Продукт удален.\n", u)

	case u.HasAction(view.ActionShowProducts):
		_, _ = b.view.ShowProductView("", u)
	}

}

func (b *BotApp) HandleAddProduct(u *tgbot.Update) {

	if u.HasAction(view.ActionAddProduct) {
		u.StartChain(string(view.ActionAddProduct)).StartChainStep("NAME").FlushChatInfo()
		_, _ = b.view.AddProductName(u)
		return
	} else if !u.IsPlainText() {
		u.FinishChain().FlushChatInfo()
		_, _ = b.view.ShowProductView("Произошла ошибка", u)
	}

	switch u.GetChainStep() {
	case "NAME":
		u.StartChainStep("DATE").AddChainData("productName", u.GetText()).FlushChatInfo()
		_, _ = b.view.AddProductDate("", u)

	case "DATE":
		date, err := time.Parse("02.01.2006", u.GetText())
		if err != nil {
			_, _ = b.view.AddProductDate("Неверный формат даты", u)
			return
		}

		err = b.prodProv.SaveProduct(model.Product{
			Id:             newUuid(),
			ChatId:         u.GetChatId(),
			Name:           u.GetChainData("productName"),
			ExpirationDate: date,
			CreatedAt:      time.Now(),
		})
		if err != nil {
			panic(err)
		}

		u.FinishChain().FlushChatInfo()

		_, _ = b.view.ShowProductView("Продукт добавлен \n", u)

	}
}

func newUuid() uuid.UUID {
	newUUID, _ := uuid.NewUUID()
	return newUUID
}
