package view

import (
	"gotestbot/sdk/tgbot"
)

const (
	ActionStart         = tgbot.Action("START")
	ActionAddProduct    = tgbot.Action("ADD_PRODUCT")
	ActionDeleteProduct = tgbot.Action("DELETE_PRODUCT")
	ActionShowProducts  = tgbot.Action("SHOW_PRODUCTS")
	ActionMerge         = tgbot.Action("MERGE")
)
