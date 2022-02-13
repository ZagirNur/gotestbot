package model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id             uuid.UUID `db:"id"`
	FridgeId       uuid.UUID `db:"fridge_id"`
	ChatId         int64     `db:"chat_id"`
	Name           string    `db:"name"`
	ExpirationDate time.Time `db:"expiration_date"`
	CreatedAt      time.Time `db:"created_at"`
}
