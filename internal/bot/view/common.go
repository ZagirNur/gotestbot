package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gotestbot/sdk/tgbot"
)

func (v *View) createButton(action tgbot.Action, data map[string]string) *tgbot.Button {
	id, _ := uuid.NewUUID()
	button := tgbot.Button{
		Id:     id.String(),
		Action: action,
		Data:   data,
	}
	err := v.chatProv.SaveButton(button)
	if err != nil {
		log.Fatal().Err(err).Msgf("")
	}
	return &button
}

func logIfError(send tgbotapi.Message, err error) (tgbotapi.Message, error) {
	if err != nil {
		log.Error().Err(err).Msgf("cannot send")
	}
	return send, err
}
