package main

import (
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/dao/pg"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/config"
	service_dao "gotestbot/internal/service/dao"
	"gotestbot/sdk/tgbot"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {

	config.InitConfig()

	if config.Conf.Dry {
		log.Info().Msg("Started in dry mode ok\nBye!")
		os.Exit(0)
	}

	InitLogger()

	pgDb := PgConnInit()
	db := initYdb()
	pgRepository := pg.NewRep(pgDb)

	bot, err := tgbot.NewBot(config.Conf.TgToken, pgRepository)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start app")
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
