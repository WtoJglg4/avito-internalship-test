package service

import (
	"github/avito/entities"
	"github/avito/pkg/repository"
)

type Authorization interface {
	CreateUser(entities.User) (int, error)
	GenerateToken(login, password string) (string, error)
}

type Banner interface {
}

type Service struct {
	Authorization
	Banner
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
	}
}
