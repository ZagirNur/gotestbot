package bot_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"gotestbot/internal/bot/dao"
	"gotestbot/internal/bot/sdk"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/service/model"
	"time"
)

type UserProvider interface {
	SaveUser(user model.User) error
	DeleteUser(userId int64) error
	GetUser(userId int64) (model.User, error)
}

type ProductProvider interface {
	GetProductsByUserId(userId int64) ([]model.Product, error)
	SaveProduct(product model.Product) error
	DeleteProduct(productId string) error
}

type BotApp struct {
	view     *view.View
	userProv UserProvider
	prodProv ProductProvider
	repos    *dao.BotRepository
}

func NewBotApp(view *view.View, userProv UserProvider, prodProv ProductProvider, repos *dao.BotRepository) *BotApp {
	return &BotApp{view: view, userProv: userProv, prodProv: prodProv, repos: repos}
}

func (b *BotApp) Handle(update tgbotapi.Update) error {
	u := sdk.NewUpdate(update, b.repos, b.repos)

	switch {
	case u.HasCommand("/start") || u.HasAction(view.ActionStart):
		u.FinishChain().FlushState()
		_, _ = b.view.StartView(u)

	case u.HasActionOrChain(view.ActionAddProduct):
		b.HandleAddProduct(u)

	case u.HasAction(view.ActionDeleteProduct):
		err := b.prodProv.DeleteProduct(u.GetButton().GetData("productId"))
		if err != nil {
			_, _ = b.view.ShowProductView("Ошибка удаления продукта.\n", u)
			return err
		}
		_, _ = b.view.ShowProductView("Продукт удален.\n", u)

	case u.HasAction(view.ActionShowProducts):
		_, _ = b.view.ShowProductView("", u)
	}

	return nil
}

func (b *BotApp) HandleAddProduct(u *sdk.Update) {

	if u.HasAction(view.ActionAddProduct) {
		u.StartChain(string(view.ActionAddProduct)).StartChainStep("NAME").FlushState()
		_, _ = b.view.AddProductName(u)
		return
	} else if !u.IsPlainText() {
		u.FinishChain().FlushState()
		_, _ = b.view.ShowProductView("Произошла ошибка", u)
	}

	switch u.GetChatState().GetChainStep() {
	case "NAME":
		u.StartChainStep("DATE").AddChainData("productName", u.GetText()).FlushState()
		_, _ = b.view.AddProductDate("", u)
	case "DATE":

		date, err := time.Parse("02.01.2006", u.GetText())
		if err != nil {
			_, _ = b.view.AddProductDate("Неверный формат даты", u)
			return
		}

		u.FinishChain().FlushState()
		err = b.prodProv.SaveProduct(model.Product{
			Id:             newUuid(),
			UserId:         u.GetChatId(),
			Name:           u.GetChatState().GetData("productName"),
			ExpirationDate: date,
			CreatedAt:      time.Now(),
		})
		if err != nil {
			panic(err)
		}

		_, _ = b.view.ShowProductView("Продукт добавлен \n", u)

	}
}

func newUuid() string {
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}

//func (b BotApp) HandleRegistration(u *sdk.Update) {
//	chatState := u.GetChatState()
//	switch chatState.GetChainStep() {
//	case "":
//		chatState.StartChain("REGISTRATION")
//		chatState.StartChainStep("NAME")
//		u.FlushState()
//		go b.view.RegistrationEnterNameView(u)
//	case "NAME":
//		if u.IsPlainText() {
//			chatState.AddData("name", u.GetText())
//			chatState.StartChainStep("AGE")
//			u.FlushState()
//			go b.view.RegistrationEnterAgeView(u, u.GetText())
//		}
//	case "AGE":
//		if u.IsPlainText() {
//			name := chatState.GetData("name")
//			age, err := strconv.Atoi(u.GetText())
//			if err != nil {
//				b.view.RegistrationEnterAgeErrorView(u, name)
//				return
//			}
//			b.userProv.SaveUser(model.User{Id: u.GetChatId(), Name: name, Age: age})
//			chatState.FinishChain()
//			u.FlushState()
//			go b.view.StartView(u)
//		}
//	}
//
//}
