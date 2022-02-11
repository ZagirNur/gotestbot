package dao

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gotestbot/sdk/tgbot"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetChat(chatId int64) (chat tgbot.ChatInfo, err error) {
	row := r.db.QueryRowx("select * from chat_info where chat_id = $1", chatId)

	if err = row.StructScan(&chat); err != nil {
		return tgbot.ChatInfo{}, errors.Wrapf(err, "unable to get chatInfo, chatId: %d", chatId)
	}
	return
}

func (r *Repository) SaveChatInfo(chat tgbot.ChatInfo) error {

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

func (r *Repository) GetButton(btnId string) (btn tgbot.Button, err error) {
	row := r.db.QueryRowx("select * from button where id = $1", btnId)

	if err = row.StructScan(&btn); err != nil {
		return tgbot.Button{}, errors.Wrapf(err, "unable to get button, btnId: %s", btnId)
	}
	return
}

func (r *Repository) SaveButton(button tgbot.Button) error {
	insert := "insert into button (id, action, data) values (:id, :action, :data)"

	if _, err := r.db.NamedExec(insert, button); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveUser(user tgbot.User) error {
	insert := `insert into profile (user_id, user_name, display_name) values (:user_id, :user_name, :display_name)
					on conflict (user_id) do update set user_name      = :user_name,
											       display_name = :display_name`

	if _, err := r.db.NamedExec(insert, user); err != nil {
		return err
	}
	return nil
}
