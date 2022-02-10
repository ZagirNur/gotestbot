package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //for db migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/config"
	"gotestbot/internal/dao/pg"
	"gotestbot/sdk/tgbot"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	config.InitConfig()

	if config.Conf.Dry {
		log.Info().Msg("Started in dry mode ok\nBye!")
		os.Exit(0)
	}

	InitLogger()

	pgDb := PgConnInit()
	pgRepository := pg.NewRepository(pgDb)

	bot, err := tgbot.NewBot(config.Conf.TgToken, pgRepository)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start app")
	}

	viewSender := view.NewView(pgRepository, pgRepository, bot)
	application := bot_handler.NewBotApp(viewSender, pgRepository)

	go func() {
		err = bot.StartLongPolling(application.Handle)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to start app")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func PgConnInit() *sqlx.DB {

	dsn := config.GetPgDsn()

	if err := MigrateDB(dsn); err != nil {
		log.Fatal().Msgf("Database migration failed: %s", err.Error())
	}
	log.Info().Msg("Database migration succeeded")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to db. dsn='%s': %s", config.DsnMaskPass(dsn), err.Error())
	}
	db.SetMaxOpenConns(config.Conf.PgMaxOpenConn)
	db.SetMaxIdleConns(config.Conf.PgMaxIdleConn)
	db.SetConnMaxLifetime(config.Conf.PgMaxLifeTime)
	db.SetConnMaxIdleTime(config.Conf.PgMaxIdleTime)
	log.Info().Msg("Connected to db")
	return db
}

func MigrateDB(dsn string) error {
	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func InitLogger() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	zerolog.TimestampFieldName = "@timestamp"

	logLvl, err := zerolog.ParseLevel(strings.ToLower(config.Conf.LogLevel))
	if err != nil {
		log.Fatal().Msgf("Failed to parse log level '%s': %s", config.Conf.LogLevel, err.Error())
	}

	zerolog.SetGlobalLevel(logLvl)

	switch config.Conf.LogFormat {
	case "plain":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	case "logstash":
		// do nothing
	default:
		log.Fatal().Msgf("Unknown log format '%s'", config.Conf.LogFormat)
	}

	log.Info().Msg("Logger successfully initialized")
}
