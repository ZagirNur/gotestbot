package dao

import (
	"context"
	"encoding/json"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"gotestbot/internal/bot/sdk"
)

var (
	roTX = table.TxControl(
		table.BeginTx(table.WithOnlineReadOnly()),
		table.CommitTx(),
	)
	rwTX = table.TxControl(
		table.BeginTx(table.WithSerializableReadWrite()),
		table.CommitTx(),
	)
)

type BotRepository struct {
	ydb table.Client
}

func NewBotRepository(ydb table.Client) *BotRepository {
	return &BotRepository{
		ydb: ydb,
	}

}

func (r *BotRepository) GetChatState(chatId int64) (sdk.ChatState, error) {

	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return sdk.ChatState{}, err
	}

	stmt, err := session.Prepare(ctx, `
			DECLARE $chat_id AS Int64?;
			SELECT * FROM chat_state WHERE chat_id = $chat_id`)
	if err != nil {
		return sdk.ChatState{}, err
	}

	idParam := table.ValueParam("$chat_id", types.OptionalValue(types.Int64Value(chatId)))
	_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
	if err != nil {
		return sdk.ChatState{}, err
	}

	res.NextResultSet(ctx, "chat_id", "active_chain", "active_chain_step", "data")
	res.NextRow()

	cs := sdk.ChatState{}
	var data []byte
	err = res.ScanWithDefaults(&cs.ChatId, &cs.ActiveChain, &cs.ActiveChainStep, &data)
	if err != nil {
		return sdk.ChatState{}, err
	}

	if err := json.Unmarshal(data, &cs.Data); err != nil {
		return sdk.ChatState{}, err
	}

	return cs, nil
}

func (r *BotRepository) SaveChatState(state sdk.ChatState) error {
	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close(context.Background()) }()

	const insert = `
			DECLARE $chat_id 			AS Int64?;
			DECLARE $active_chain 		AS Utf8?;
			DECLARE $active_chain_step 	AS Utf8?;
			DECLARE $data 				AS Utf8?;

			UPSERT INTO chat_state (chat_id, active_chain, active_chain_step, data)
							VALUES ($chat_id, $active_chain, $active_chain_step, $data);`

	stmt, err := session.Prepare(ctx, insert)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(state.Data)
	_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
		table.ValueParam("$chat_id", types.OptionalValue(types.Int64Value(state.ChatId))),
		table.ValueParam("$active_chain", types.OptionalValue(types.UTF8Value(string(state.ActiveChain)))),
		table.ValueParam("$active_chain_step", types.OptionalValue(types.UTF8Value(state.ActiveChainStep))),
		table.ValueParam("$data", types.OptionalValue(types.UTF8Value(string(data)))),
	))

	if err != nil {
		return err
	}
	return nil

}

func (r *BotRepository) GetButton(buttonId string) (sdk.Button, error) {

	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return sdk.Button{}, err
	}

	stmt, err := session.Prepare(ctx, `
			DECLARE $id AS Utf8?;
			SELECT * FROM button WHERE id = $id`)
	if err != nil {
		return sdk.Button{}, err
	}

	idParam := table.ValueParam("$id", types.OptionalValue(types.UTF8Value(buttonId)))
	_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
	if err != nil {
		return sdk.Button{}, err
	}

	res.NextResultSet(ctx, "id", "action", "data")
	res.NextRow()

	b := sdk.Button{}
	var data []byte
	var action = ""
	err = res.ScanWithDefaults(&b.Id, &action, &data)
	if err != nil {
		return sdk.Button{}, err
	}

	b.Action = sdk.Action(action)
	if err := json.Unmarshal(data, &b.Data); err != nil {
		return sdk.Button{}, err
	}

	return b, nil
}

func (r *BotRepository) SaveButton(button sdk.Button) error {
	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close(context.Background()) }()

	const insert = `
			DECLARE $id AS Utf8?;
			DECLARE $action AS Utf8?;
			DECLARE $data AS Utf8?;
			UPSERT INTO button
				(id, action, data)
			VALUES
				($id, $action, $data);`

	stmt, err := session.Prepare(ctx, insert)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(button.Data)
	_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
		table.ValueParam("$id", types.OptionalValue(types.UTF8Value(button.Id))),
		table.ValueParam("$action", types.OptionalValue(types.UTF8Value(string(button.Action)))),
		table.ValueParam("$data", types.OptionalValue(types.UTF8Value(string(data)))),
	))

	if err != nil {
		return err
	}
	return nil

}
