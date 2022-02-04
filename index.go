package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gotestbot/internal"
	"gotestbot/internal/view"
	"io"
	"net/http"
	"os"
)

func Handler(rw http.ResponseWriter, req *http.Request) {
	fmt.Println(req)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	update, err := bot.HandleUpdate(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(update)

	repository := internal.NewRepository(initYdb())
	viewSender := view.NewView(repository, repository, bot)

	err = internal.NewBotApp(viewSender, repository).Handle(*update)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("X-Custom-Header", "Test")
	rw.WriteHeader(200)
	name := req.URL.Query().Get("name")
	io.WriteString(rw, fmt.Sprintf("Hello, %s!", name))
}
