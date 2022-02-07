package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/dao"
	"gotestbot/internal/bot/view"
	service_dao "gotestbot/internal/service/dao"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req)

	tgApi, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	update, err := tgApi.HandleUpdate(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(update)

	client := initYdb()
	rep := dao.NewBotRepository(client)
	serviceRep := service_dao.NewRepository(client)
	viewSender := view.NewView(rep, serviceRep, serviceRep, tgApi)

	application := bot_handler.NewBotApp(viewSender, serviceRep, serviceRep, rep)
	err = application.Handle(*update)
	if err != nil {
		panic(err)
	}

	rw.WriteHeader(200)
}
