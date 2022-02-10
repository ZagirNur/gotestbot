package model

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ChatId         int64     `db:"chat_id"`
	Id             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	ExpirationDate time.Time `db:"expiration_date"`
	CreatedAt      time.Time `db:"created_at"`
}
