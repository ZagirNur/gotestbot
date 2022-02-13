package dao

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gotestbot/internal/service/model"
)

func (r *Repository) SaveProduct(product model.Product) error {
	insert := `INSERT INTO product (id, fridge_id, name, expiration_date, created_at) 
		VALUES (:id, (SELECT fridge_id FROM chat_fridge WHERE chat_id = :chat_id), :name, :expiration_date, :created_at)`

	if _, err := r.db.NamedExec(insert, product); err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteProduct(productId string) error {
	_, err := r.db.Exec(`DELETE FROM product WHERE id=$1`, productId)
	if err != nil {
		return errors.Wrapf(err, "unable to delte product, productId: %s", productId)
	}
	return nil
}

func (r *Repository) GetProductsByChatId(chatId int64) (products []model.Product, err error) {
	rows, err := r.db.Queryx(`SELECT product.* FROM product 
    							JOIN fridge f ON f.id = product.fridge_id 
    							JOIN chat_fridge cf ON f.id = cf.fridge_id 
								WHERE chat_id = $1`, chatId)
	defer rows.Close()

	for rows.Next() {
		p := model.Product{}
		if err = rows.StructScan(&p); err != nil {
			return []model.Product{}, errors.Wrapf(err, "unable to get products, chatId: %d", chatId)
		}
		p.ChatId = chatId

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) ExistsFridgeByChatId(chatId int64) (exists bool, err error) {
	row := r.db.QueryRow("SELECT exists(SELECT 1 FROM chat_fridge WHERE chat_id = $1)", chatId)
	err = row.Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *Repository) CreateFridge() (id uuid.UUID, err error) {
	id = uuid.New()

	_, err = r.db.Exec("INSERT INTO fridge (id) VALUES ($1)", id.String())
	if err != nil {
		return uuid.UUID{}, nil
	}
	return id, nil
}

func (r *Repository) GetFridgeByChatId(chatId int64) (id uuid.UUID, err error) {
	row := r.db.QueryRow("SELECT fridge_id FROM chat_fridge WHERE chat_id = $1", chatId)
	err = row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r *Repository) AddFridgeToChat(chatId int64, fridgeId uuid.UUID) error {
	_, err := r.db.Exec(`INSERT INTO chat_fridge (chat_id, fridge_id) VALUES ($1, $2)
						ON CONFLICT (chat_id) DO UPDATE SET fridge_id = $2;`, chatId, fridgeId.String())
	if err != nil {
		return err
	}
	return nil
}
