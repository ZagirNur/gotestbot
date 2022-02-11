package dao

import (
	"github.com/pkg/errors"
	"gotestbot/internal/service/model"
)

func (r *Repository) SaveProduct(product model.Product) error {
	insert := `insert into product (chat_id, id, name, expiration_date, created_at) 
							values (:chat_id, :id, :name, :expiration_date, :created_at)`

	if _, err := r.db.NamedExec(insert, product); err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteProduct(productId string) error {
	_, err := r.db.Exec(`delete from product where id=$1`, productId)
	if err != nil {
		return errors.Wrapf(err, "unable to delte product, productId: %s", productId)
	}
	return nil
}

func (r *Repository) GetProductsByChatId(chatId int64) (products []model.Product, err error) {
	rows, err := r.db.Queryx("select * from product where chat_id = $1", chatId)
	defer rows.Close()

	for rows.Next() {
		p := model.Product{}
		if err = rows.StructScan(&p); err != nil {
			return []model.Product{}, errors.Wrapf(err, "unable to get products, chatId: %d", chatId)
		}
		products = append(products, p)

	}
	return products, nil
}
