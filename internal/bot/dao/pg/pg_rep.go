package pg

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"gotestbot/internal/bot/sdk"
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

func (r Rep) GetButton(btnId string) (btn sdk.Button, err error) {
	row := r.db.QueryRowx("select * from button where id = $1", btnId)

	if err = row.StructScan(&btn); err != nil {
		return sdk.Button{}, err
	}
	return
}

func (r Rep) SaveButton(button sdk.Button) error {
	insert := "insert into button (id, action, data) values (:id, :action, :data)"

	if _, err := r.db.NamedExec(insert, button); err != nil {
		return err
	}
	return nil
}
