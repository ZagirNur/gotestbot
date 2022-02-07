package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc"
	"gotestbot/internal/bot/bot_handler"
	"gotestbot/internal/bot/dao"
	"gotestbot/internal/bot/view"
	service_dao "gotestbot/internal/service/dao"
	"os"
	"path"
	"time"
)

func main() {
	api, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		return
	}

	client := initYdb()
	rep := dao.NewBotRepository(client)
	serviceRep := service_dao.NewRepository(client)
	viewSender := view.NewView(rep, serviceRep, serviceRep, api)

	application := bot_handler.NewBotApp(viewSender, serviceRep, serviceRep, rep)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range api.GetUpdatesChan(u) {
		fmt.Println(update)
		err = application.Handle(update)
		if err != nil {
			panic(err)
		}
	}
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
