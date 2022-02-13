package bot_handler

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/service"
	"gotestbot/internal/service/model"
	"gotestbot/sdk/tgbot"
)

type ProductProvider interface {
	GetProductsByChatId(chatId int64) ([]model.Product, error)
	SaveProduct(product model.Product) error
	DeleteProduct(productId string) error
	CreateMergedFridge(...string) (uuid.UUID, error)
	SetFridge(fridge uuid.UUID, chatIds ...int64)
}

type BotApp struct {
	view        *view.View
	prodService *service.ProdService
}

func NewBotApp(view *view.View, prodProv *service.ProdService) *BotApp {
	return &BotApp{view: view, prodService: prodProv}
}

func (b *BotApp) Handle(u *tgbot.Update) {

	if err := b.prodService.CreateFridgeIfNotExists(u.GetChatId()); err != nil {
		log.Error().Err(err).Msgf("cannot CreateFridgeIfNotExists, chatId=%d", u.GetChatId())
	}

	switch {
	case u.HasCommand("/start") || u.HasAction(view.ActionStart):
		u.FinishChain().FlushChatInfo()
		_, _ = b.view.StartView(u)

	case u.HasCommand("/share") || u.GetInline() == "share" || u.HasAction(view.ActionMerge):
		u.FinishChain().FlushChatInfo()
		b.handleShare(u)

	case u.HasActionOrChain(view.ActionAddProduct):
		b.HandleAddProduct(u)

	case u.HasAction(view.ActionDeleteProduct):
		err := b.prodService.DeleteProduct(u.GetButton().GetData("productId"))
		if err != nil {
			_, _ = b.view.ShowProductView("Ошибка удаления продукта.\n", u)
			return
		}
		_, _ = b.view.ShowProductView("Продукт удален.\n", u)

	case u.HasAction(view.ActionShowProducts):
		_, _ = b.view.ShowProductView("", u)
	}

}
