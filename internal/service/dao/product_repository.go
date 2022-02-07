package dao

import (
	"context"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"gotestbot/internal/service/model"
)

func (r *Repository) GetProductsByUserId(userId int64) (products []model.Product, err error) {
	ctx := context.Background()
	return products, r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

		stmt, err := session.Prepare(ctx, `
			DECLARE $user_id AS Int64?;
			SELECT * FROM product WHERE user_id = $user_id`)
		if err != nil {
			return err
		}

		idParam := table.ValueParam("$user_id", types.OptionalValue(types.Int64Value(userId)))
		_, res, err := stmt.Execute(ctx, roTX, table.NewQueryParameters(idParam))
		if err != nil {
			return err
		}

		res.NextResultSet(ctx, "user_id", "id", "name", "expiration_date", "created_at")

		for res.NextRow() {

			p := model.Product{}
			err = res.ScanWithDefaults(&p.UserId, &p.Id, &p.Name, &p.ExpirationDate, &p.CreatedAt)
			if err != nil {
				return err
			}
			products = append(products, p)
		}

		return nil
	})
}

func (r *Repository) SaveProduct(p model.Product) error {
	ctx := context.Background()
	return r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

		const insert = `
			DECLARE $user_id			AS Int64?;
			DECLARE $id					AS Utf8?;
			DECLARE $name				AS Utf8?;
			DECLARE $expiration_date	AS Date?;
			DECLARE $created_at			AS DateTime?;

			UPSERT INTO product (user_id, id, name, expiration_date, created_at)
						VALUES 	($user_id, $id, $name, $expiration_date, $created_at);`

		stmt, err := session.Prepare(ctx, insert)
		if err != nil {
			return err
		}

		_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
			table.ValueParam("$user_id", types.OptionalValue(types.Int64Value(p.UserId))),
			table.ValueParam("$id", types.OptionalValue(types.UTF8Value(p.Id))),
			table.ValueParam("$name", types.OptionalValue(types.UTF8Value(p.Name))),
			table.ValueParam("$expiration_date", types.OptionalValue(types.DateValueFromTime(p.ExpirationDate))),
			table.ValueParam("$created_at", types.OptionalValue(types.DatetimeValueFromTime(p.CreatedAt))),
		))
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *Repository) DeleteProduct(productId string) error {
	ctx := context.Background()
	return r.ydb.Do(ctx, func(ctx context.Context, session table.Session) error {

		const insert = `
			DECLARE $id		AS Utf8?;
			DELETE FROM product WHERE id = $id;`

		stmt, err := session.Prepare(ctx, insert)
		if err != nil {
			return err
		}

		_, _, err = stmt.Execute(ctx, rwTX, table.NewQueryParameters(
			table.ValueParam("$id", types.OptionalValue(types.UTF8Value(productId))),
		))
		if err != nil {
			return err
		}
		return nil
	})
}
