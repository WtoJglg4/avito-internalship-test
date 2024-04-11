package service

import (
	"github/avito/entities"
	"github/avito/pkg/repository"
)

type Authorization interface {
	CreateUser(entities.User) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, string, error)
}

type Banner interface {
	CreateBanner(entities.Banner) (int, error)
	GetAllBanners(entities.QueryFilters) ([]entities.Banner, error)
	DeleteBanners(entities.QueryFilters) error
	UserBanner(entities.QueryFilters) (entities.Content, error)
}

type Service struct {
	Authorization
	Banner
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		Banner:        NewBannersService(repo),
	}
}
