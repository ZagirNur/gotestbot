package main

import (
	"context"
	"crypto/tls"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"github.com/yandex-cloud/ydb-go-sdk"
	"github.com/yandex-cloud/ydb-go-sdk/auth/iam"
	"github.com/yandex-cloud/ydb-go-sdk/table"
	"gotestbot/internal"
	"gotestbot/internal/view"
	"os"
)

func main() {
	api, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		return
	}

	client := initYdb()
	repository := internal.NewRepository(client)
	viewSender := view.NewView(repository, repository, api)

	application := internal.NewBotApp(viewSender, repository)

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

func initYdb() *table.Client {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	credentials, err := iam.NewClient(
		iam.WithServiceFile("./SA-FILE"),
		iam.WithDefaultEndpoint(),
		iam.WithSystemCertPool(),
	)
	if err != nil {
		log.Fatal().Err(err)
	}

	dialer := &ydb.Dialer{
		DriverConfig: &ydb.DriverConfig{
			Database:    os.Getenv("YDB_DATABASE"),
			Credentials: credentials,
		},
		TLSConfig: &tls.Config{},
	}

	driver, err := dialer.Dial(ctx, "ydb.serverless.yandexcloud.net:2135")
	if err != nil {
		log.Fatal().Err(err)
	}

	return &table.Client{
		Driver: driver,
	}
}
