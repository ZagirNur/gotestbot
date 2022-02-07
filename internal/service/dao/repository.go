package dao

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"gotestbot/internal/service/model"
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

type Repository struct {
	ydb table.Client
}

func NewRepository(ydb table.Client) *Repository {
	return &Repository{ydb: ydb}
}

func (r *Repository) GetUser(userId int64) (model.User, error) {

	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return model.User{}, err
	}

	stmt, err := session.Prepare(ctx, `
			DECLARE $id AS Int64?;
			SELECT * FROM users WHERE id = $id`)
	if err != nil {
		return model.User{}, err
	}

	idParam := table.ValueParam("$id", types.OptionalValue(types.Int64Value(userId)))
	_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
	if err != nil {
		return model.User{}, err
	}

	res.NextResultSet(ctx, "id", "age", "name")
	res.NextRow()

	cs := model.User{}
	err = res.ScanWithDefaults(&cs.Id, &cs.Age, &cs.Name)
	if err != nil {
		return model.User{}, err
	}

	return cs, nil
}

func (r *Repository) SaveUser(user model.User) error {
	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close(context.Background()) }()

	const insert = `
			DECLARE $id		AS Int64?;
			DECLARE $age	AS Int32?;
			DECLARE $name	AS Utf8?;

			UPSERT INTO users (id, age, name)
							VALUES ($id, $age, $name);`

	stmt, err := session.Prepare(ctx, insert)
	if err != nil {
		return err
	}

	_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
		table.ValueParam("$id", types.OptionalValue(types.Int64Value(user.Id))),
		table.ValueParam("$age", types.OptionalValue(types.Int32Value(int32(user.Age)))),
		table.ValueParam("$name", types.OptionalValue(types.UTF8Value(user.Name))),
	))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteUser(userId int64) error {
	ctx := context.Background()
	session, err := r.ydb.CreateSession(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = session.Close(context.Background()) }()

	const insert = `
			DECLARE $id		AS Int64?;

			DELETE FROM users WHERE id = $id;`

	stmt, err := session.Prepare(ctx, insert)
	if err != nil {
		return err
	}

	_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
		table.ValueParam("$id", types.OptionalValue(types.Int64Value(userId))),
	))
	if err != nil {
		return err
	}
	return nil
}