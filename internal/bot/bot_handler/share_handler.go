package bot_handler

import (
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/view"
	"gotestbot/sdk/tgbot"
	"strconv"
)

func (b *BotApp) handleShare(u *tgbot.Update) {

	switch {
	case u.HasCommand("/share"):
		_, _ = b.view.Share(u)

	case u.GetInline() == "share":
		_, _ = b.view.ShareInline(u)

	case u.HasAction(view.ActionMerge):
		sChatId := u.GetButton().GetData("chatId")
		chatId, err := strconv.ParseInt(sChatId, 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("cannot merge fridges")
			_, _ = b.view.ErrorMessage(u, "Ошибка объединения холодильников")
			return
		}
		fridgeId, err := b.prodService.GetFridgeByChatId(chatId)
		if err != nil {
			log.Error().Err(err).Msg("cannot merge fridges")
			_, _ = b.view.ErrorMessage(u, "Ошибка объединения холодильников")
			return
		}

		if err := b.prodService.AddFridgeToChat(u.GetChatId(), fridgeId); err != nil {
			log.Error().Err(err).Msg("cannot merge fridges")
			_, _ = b.view.ErrorMessage(u, "Ошибка объединения холодильников")
			return
		}

		_, _ = b.view.GoToBotScreen(u)
	}

}
