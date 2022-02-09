package main

import (
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/dao/pg"
	"gotestbot/internal/bot/view"
	service_dao "gotestbot/internal/service/dao"
	"gotestbot/sdk/tgbot"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	tgToken := os.Getenv("TG_TOKEN")

	db := initYdb()
	pgRepository := pg.NewRep()

	bot, err := tgbot.NewBot(tgToken, pgRepository)
	if err != nil {
		log.Error().Err(err).Msg("unable to start app")
		return
	}

	serviceRep := service_dao.NewRepository(db)
	viewSender := view.NewView(pgRepository, serviceRep, bot)
	application := bot_handler.NewBotApp(viewSender, serviceRep, serviceRep)

	update, err := bot.WrapRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("unable read request")
		return
	}

	application.Handle(update)

	rw.WriteHeader(200)
}
