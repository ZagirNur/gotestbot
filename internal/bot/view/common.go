package view

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/sdk"
)

func (v *View) CreateButton(action sdk.Action, data map[string]string) *sdk.Button {

	id, _ := uuid.NewUUID()
	button := sdk.Button{
		Id:     id.String(),
		Action: action,
		Data:   data,
	}
	//go func() {
	err := v.btnProv.SaveButton(button)
	if err != nil {
		log.Fatal().Err(err).Msgf("")
	}
	//}()
	return &button
}

func logIfError(send tgbotapi.Message, err error) (tgbotapi.Message, error) {
	if err != nil {
		log.Error().Err(err).Msgf("cannot send")
	}
	return send, err
}
