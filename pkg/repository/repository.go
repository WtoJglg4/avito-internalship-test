package repository

import (
	"github/avito/entities"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	GetUser(login, password string) (entities.User, error)
}

type Banner interface {
	CreateBanner(entities.Banner) (int, error)
	GetAllBanners(entities.QueryFilters) ([]entities.Banner, error)
	DeleteBanners(entities.QueryFilters) error
	UserBanner(entities.QueryFilters) (entities.Content, error)
}

type Repository struct {
	Authorization
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banner:        NewBannersPostgres(db),
	}
}
