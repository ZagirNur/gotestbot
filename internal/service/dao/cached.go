package dao

import (
	"gotestbot/internal/bot/dao"
	"gotestbot/internal/bot/sdk"
	"sync"
)

type Cached struct {
	*dao.BotRepository

	buttons map[string]sdk.Button
	chats   map[int64]sdk.ChatState
	m       sync.Mutex
}

func (c *Cached) GetButton(buttonId string) (sdk.Button, error) {
	c.m.Lock()
	defer c.m.Unlock()
	if b, ok := c.buttons[buttonId]; ok {
		return b, nil
	}
	return c.BotRepository.GetButton(buttonId)
}

func (c *Cached) SaveButton(button sdk.Button) error {
	c.m.Lock()
	defer c.m.Unlock()
	c.buttons[button.Id] = button
	return nil
}

func (c *Cached) FlushCache(button sdk.Button) error {
	c.m.Lock()
	defer c.m.Unlock()
	//todo c.BotRepository.SaveAllButtons()
	return nil
}
