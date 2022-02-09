package pg

import (
	"encoding/json"
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
	row := r.db.QueryRow("select * from button where id = $1", btnId)
	var data []byte
	dataMap := map[string]string{}
	err = row.Scan(&btn.Id, &btn.Action, &data)
	if err != nil {
		return sdk.Button{}, err
	}
	json.Unmarshal(data, &dataMap)
	btn.Data = dataMap
	return
}

func (r Rep) SaveButton(button sdk.Button) error {
	data, _ := json.Marshal(button.Data)

	_, err := r.db.Exec("insert into button (id, action, data) values ($1,$2,$3::json)",
		button.Id, button.Action, data)
	if err != nil {
		return err
	}
	return nil
}
