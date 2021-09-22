package store

import "bot/models"

type UserRepository interface {
	Create(*models.User) error
	FindByUsername(username string) (*models.User, error)
	FindByChatId(chatId int) (*models.User, error)
	GetAll() ([]models.User, error)
}