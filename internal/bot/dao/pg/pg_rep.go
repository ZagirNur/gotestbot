package pg

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gotestbot/sdk/tgbot"
	"net/url"
)

type Rep struct {
	db *sqlx.DB
}

func NewRep() *Rep {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		url.QueryEscape("postgres"),
		url.QueryEscape("postgres"),
		"localhost",
		5432,
		"gotestbot")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to connect to db.")
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	log.Info().Msg("Connected to db")

	return &Rep{db: db}
}

func (r *Rep) GetChat(chatId int64) (chat tgbot.ChatInfo, err error) {
	row := r.db.QueryRowx("select * from chat_info where chat_id = $1", chatId)

	if err = row.StructScan(&chat); err != nil {
		return tgbot.ChatInfo{}, errors.Wrapf(err, "unable to get chatInfo, chatId: %d", chatId)
	}
	return
}

func (r *Rep) SaveChatInfo(chat tgbot.ChatInfo) error {

	insert := `insert into chat_info (chat_id, active_chain, active_chain_step, chain_data)
								values (:chat_id, :active_chain, :active_chain_step, :chain_data)
								on conflict (chat_id) do update set active_chain      = :active_chain,
														  active_chain_step = :active_chain_step,
														  chain_data        = :chain_data`

	if _, err := r.db.NamedExec(insert, chat); err != nil {
		return errors.Wrap(err, "unable to save chatInfo")
	}
	return nil
}

func (r *Rep) GetButton(btnId string) (btn tgbot.Button, err error) {
	row := r.db.QueryRowx("select * from button where id = $1", btnId)

	if err = row.StructScan(&btn); err != nil {
		return tgbot.Button{}, errors.Wrapf(err, "unable to get button, btnId: %s", btnId)
	}
	return
}

func (r *Rep) SaveButton(button tgbot.Button) error {
	insert := "insert into button (id, action, data) values (:id, :action, :data)"

	if _, err := r.db.NamedExec(insert, button); err != nil {
		return err
	}
	return nil
}
