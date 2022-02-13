package bot_handler

import (
	"github.com/google/uuid"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/service/model"
	"gotestbot/sdk/tgbot"
	"time"
)

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
			_, _ = b.view.AddProductDate("Неверный формат даты\n", u)
			return
		}

		err = b.prodService.SaveProduct(model.Product{
			Id:             uuid.New(),
			ChatId:         u.GetChatId(),
			Name:           u.GetChainData("productName"),
			ExpirationDate: date,
			CreatedAt:      time.Now(),
		})
		if err != nil {
			if u.IsButton() {
				_, _ = b.view.ErrorMessage(u, "Не удалось сохранить продукт\n")
			} else {
				_, _ = b.view.AddProductDate("Не удалось сохранить продукт\n", u)
			}
			return
		}

		u.FinishChain().FlushChatInfo()

		_, _ = b.view.ShowProductView("Продукт добавлен \n", u)

	}
}
