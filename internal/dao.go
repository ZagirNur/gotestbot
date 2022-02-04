package internal

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/yandex-cloud/ydb-go-sdk"
	"github.com/yandex-cloud/ydb-go-sdk/table"
	"github.com/yandex-cloud/ydb-go-sdk/ydbsql"
	"gotestbot/sdk"
)

type Repository struct {
	ydb *table.Client

	users   map[int64]sdk.User
	buttons map[string]sdk.Button
	states  map[int64]sdk.ChatState
}

func NewRepository(ydb *table.Client) *Repository {
	return &Repository{
		ydb:     ydb,
		users:   map[int64]sdk.User{},
		buttons: map[string]sdk.Button{},
		states:  map[int64]sdk.ChatState{},
	}

}

func (r *Repository) GetChatState(chatId int64) sdk.ChatState {

	state, ok := r.states[chatId]
	if ok {
		return state
	}
	return sdk.ChatState{
		ChatId: chatId,
	}
}

func (r *Repository) SaveChatState(state sdk.ChatState) {
	r.states[state.ChatId] = state
}

func (r *Repository) GetButton(buttonId string) sdk.Button {
	return r.buttons[buttonId]
}

func (r *Repository) SaveUser(user sdk.User) {
	r.users[user.Id] = user
}

func (r *Repository) DeleteUser(userId int64) {
	delete(r.users, userId)
}

func (r *Repository) GetUser(userId int64) sdk.User {
	return r.users[userId]
}

func (r *Repository) SaveButton(button sdk.Button) sdk.Button {
	r.buttons[button.Id] = button
	return button
}

func (r *Repository) saveMsg(msg string) error {

	session, err := r.ydb.CreateSession(context.Background())
	if err != nil {
		panic(err)
	}

	parameters := table.NewQueryParameters()
	parameters.Add()
	random, _ := uuid.NewRandom()
	control := table.TxControl(
		table.BeginTx(
			table.WithSerializableReadWrite(),
		),
		table.CommitTx(),
	)
	_, _, err = session.Execute(context.Background(), control, ""+
		"DECLARE $id AS String; "+
		"DECLARE $msg AS String; "+
		"INSERT INTO messages(id, msg) values ($id, $msg);",
		table.NewQueryParameters(
			table.ValueParam("$id", ydb.StringValue([]byte(random.String()))),
			table.ValueParam("$msg", ydb.StringValue([]byte(msg))),
		))
	if err != nil {
		panic(err)
	}

	return nil
}

func (r *Repository) getLastMsg() (string, error) {

	db := sql.OpenDB(ydbsql.Connector(ydbsql.WithClient(r.ydb)))
	row := db.QueryRow("SELECT msg FROM messages LIMIT 1")
	s := ""
	err := row.Scan(&s)
	if err != nil {
		panic(err)
	}
	return s, nil
}
