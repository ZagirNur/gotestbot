package view

import (
	"github.com/google/uuid"
	"gotestbot/sdk"
)

func (v *View) CreateButton(action string) *sdk.Button {

	id, _ := uuid.NewUUID()
	button := sdk.Button{
		Id:     id.String(),
		Action: action,
	}
	v.btnProv.SaveButton(button)
	return &button
}
