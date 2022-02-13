package service

import (
	"github.com/pkg/errors"
	"gotestbot/internal/dao"
)

type ProdService struct {
	*dao.Repository
}

func NewProdService(repository *dao.Repository) *ProdService {
	return &ProdService{Repository: repository}
}

func (s ProdService) CreateFridgeIfNotExists(chatId int64) error {
	exists, err := s.Repository.ExistsFridgeByChatId(chatId)
	if err != nil {
		return errors.Wrapf(err, "cannot querying for fridge existing. chatId=%d", chatId)
	}
	if !exists {
		fridgeId, err := s.Repository.CreateFridge()
		if err != nil {
			return errors.Wrapf(err, "cannot create new fridge. chatId=%d, fridgeId=%s", chatId, fridgeId)
		}

		if err := s.Repository.AddFridgeToChat(chatId, fridgeId); err != nil {
			return errors.Wrapf(err, "cannot querying for fridge existing. chatId=%d", chatId)
		}
	}
	return nil

}
