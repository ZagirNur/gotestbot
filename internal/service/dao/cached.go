package dao

import (
	"gotestbot/internal/bot/dao"
	"gotestbot/sdk/tgbot"
	"sync"
)

type Cached struct {
	*dao.BotRepository

	buttons map[string]tgbot.Button
	chats   map[int64]tgbot.ChatInfo
	m       sync.Mutex
}

func (c *Cached) GetButton(buttonId string) (tgbot.Button, error) {
	c.m.Lock()
	defer c.m.Unlock()
	if b, ok := c.buttons[buttonId]; ok {
		return b, nil
	}
	return c.BotRepository.GetButton(buttonId)
}

func (c *Cached) SaveButton(button tgbot.Button) error {
	c.m.Lock()
	defer c.m.Unlock()
	c.buttons[button.Id] = button
	return nil
}

func (c *Cached) FlushCache(button tgbot.Button) error {
	c.m.Lock()
	defer c.m.Unlock()
	//todo c.BotRepository.SaveAllButtons()
	return nil
}
