package view

import "gotestbot/internal/bot/sdk"

const (
	CommandStart      = sdk.Command("/start")
	CommandNewShelf   = sdk.Command("/new_shelf")
	CommandNewProduct = sdk.Command("/new_product")
)

const (
	ActionStart         = sdk.Action("START")
	ActionAddProduct    = sdk.Action("ADD_PRODUCT")
	ActionDeleteProduct = sdk.Action("DELETE_PRODUCT")
	ActionShowProducts  = sdk.Action("SHOW_PRODUCTS")
)
