package main

import (
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/dao"
	"gotestbot/sdk/tgbot"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {

	InitConfig()

	if conf.Dry {
		log.Info().Msg("Started in dry mode ok\nBye!")
		os.Exit(0)
	}

	InitLogger()

	pgDb := PgConnInit()
	pgRepository := dao.NewRepository(pgDb)

	bot, err := tgbot.NewBot(conf.TgToken, pgRepository)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start app")
	}

	viewSender := view.NewView(pgRepository, pgRepository, bot)
	application := bot_handler.NewBotApp(viewSender, pgRepository)
	update, err := bot.WrapRequest(req)
	if err != nil {
		log.Error().Err(err).Msg("unable read request")
		return
	}

	application.Handle(update)

	rw.WriteHeader(200)
}
