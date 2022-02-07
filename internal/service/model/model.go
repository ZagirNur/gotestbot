package model

import "time"

type User struct {
	Id   int64
	Name string
	Age  int
}

type Product struct {
	UserId         int64
	Id             string
	Name           string
	ExpirationDate time.Time
	CreatedAt      time.Time
}
