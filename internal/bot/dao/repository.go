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

func (r *BotRepository) GetChatState(chatId int64) (cs sdk.ChatState, err error) {
	ctx := context.Background()

	return cs, r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

		stmt, err := session.Prepare(ctx, `
			DECLARE $chat_id AS Int64?;
			SELECT * FROM chat_state WHERE chat_id = $chat_id`)
		if err != nil {
			return err
		}

		idParam := table.ValueParam("$chat_id", types.OptionalValue(types.Int64Value(chatId)))
		_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
		if err != nil {
			return err
		}

		res.NextResultSet(ctx, "chat_id", "active_chain", "active_chain_step", "data")
		res.NextRow()

		var data []byte
		err = res.ScanWithDefaults(&cs.ChatId, &cs.ActiveChain, &cs.ActiveChainStep, &data)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, &cs.Data); err != nil {
			return err
		}

		return nil
	})
}

func (r *BotRepository) SaveChatState(state sdk.ChatState) error {
	ctx := context.Background()

	return r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

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
	})

}

func (r *BotRepository) GetButton(buttonId string) (btn sdk.Button, err error) {

	ctx := context.Background()

	return btn, r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

		stmt, err := session.Prepare(ctx, `
			DECLARE $id AS Utf8?;
			SELECT * FROM button WHERE id = $id`)
		if err != nil {
			return err
		}

		idParam := table.ValueParam("$id", types.OptionalValue(types.UTF8Value(buttonId)))
		_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
		if err != nil {
			return err
		}

		res.NextResultSet(ctx, "id", "action", "data")
		res.NextRow()

		var data []byte
		var action = ""
		err = res.ScanWithDefaults(&btn.Id, &action, &data)
		if err != nil {
			return err
		}

		btn.Action = sdk.Action(action)
		if err := json.Unmarshal(data, &btn.Data); err != nil {
			return err
		}

		return nil
	})
}

func (r *BotRepository) SaveButton(button sdk.Button) error {
	ctx := context.Background()
	return r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

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

	})
}
