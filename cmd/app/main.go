package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //for db migration
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/dao/pg"
	"gotestbot/internal/bot/view"
	"gotestbot/internal/config"
	service_dao "gotestbot/internal/service/dao"
	"gotestbot/sdk/tgbot"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

func main() {

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

func initYdb() table.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := os.Getenv("YDB_DATABASE")
	db, err := ydb.New(ctx,
		ydb.WithDialTimeout(10*time.Second),
		ydb.WithConnectionString("grpcs://ydb.serverless.yandexcloud.net:2135"),
		yc.WithServiceAccountKeyFileCredentials("./SA-FILE", yc.WithSystemCertPool()),
		ydb.WithDatabase(database),
		ydb.WithSessionPoolKeepAliveMinSize(5),
		ydb.WithSessionPoolIdleThreshold(10*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	client := db.Table()

	createTables(ctx, client, database)

	return client
}

func createTables(ctx context.Context, client table.Client, dbPrefix string) {
	err := client.Do(ctx, func(ctx context.Context, s table.Session) error {
		return s.CreateTable(ctx, path.Join(dbPrefix, "button"),
			options.WithColumn("id", types.Optional(types.TypeUTF8)),
			options.WithColumn("action", types.Optional(types.TypeUTF8)),
			options.WithColumn("data", types.Optional(types.TypeUTF8)),
			options.WithPrimaryKeyColumn("id"),
		)
	})
	if err != nil {
		panic(err)
	}

	err = client.Do(ctx, func(ctx context.Context, s table.Session) error {
		return s.CreateTable(ctx, path.Join(dbPrefix, "chat_state"),
			options.WithColumn("chat_id", types.Optional(types.TypeInt64)),
			options.WithColumn("active_chain", types.Optional(types.TypeUTF8)),
			options.WithColumn("active_chain_step", types.Optional(types.TypeUTF8)),
			options.WithColumn("data", types.Optional(types.TypeUTF8)),
			options.WithPrimaryKeyColumn("chat_id"),
		)
	})
	if err != nil {
		panic(err)
	}

	err = client.Do(ctx, func(ctx context.Context, s table.Session) error {
		return s.CreateTable(ctx, path.Join(dbPrefix, "users"),
			options.WithColumn("id", types.Optional(types.TypeInt64)),
			options.WithColumn("age", types.Optional(types.TypeInt32)),
			options.WithColumn("name", types.Optional(types.TypeUTF8)),
			options.WithPrimaryKeyColumn("id"),
		)
	})
	if err != nil {
		panic(err)
	}

	err = client.Do(ctx, func(ctx context.Context, s table.Session) error {
		return s.CreateTable(ctx, path.Join(dbPrefix, "product"),
			options.WithColumn("user_id", types.Optional(types.TypeInt64)),
			options.WithColumn("id", types.Optional(types.TypeUTF8)),
			options.WithColumn("name", types.Optional(types.TypeUTF8)),
			options.WithColumn("expiration_date", types.Optional(types.TypeDate)),
			options.WithColumn("created_at", types.Optional(types.TypeDatetime)),
			options.WithPrimaryKeyColumn("user_id", "id"),
		)
	})
	if err != nil {
		panic(err)
	}
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
