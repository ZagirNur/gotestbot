package view

import (
	"gotestbot/sdk/tgbot"
)

const (
	CommandStart      = tgbot.Command("/start")
	CommandNewShelf   = tgbot.Command("/new_shelf")
	CommandNewProduct = tgbot.Command("/new_product")
)

const (
	ActionStart         = tgbot.Action("START")
	ActionAddProduct    = tgbot.Action("ADD_PRODUCT")
	ActionDeleteProduct = tgbot.Action("DELETE_PRODUCT")
	ActionShowProducts  = tgbot.Action("SHOW_PRODUCTS")
)
