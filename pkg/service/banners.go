package service

import (
	"github/avito/entities"
	"github/avito/pkg/repository"
)

type BannersService struct {
	repo *repository.Repository
}

func NewBannersService(repo *repository.Repository) *BannersService {
	return &BannersService{repo: repo}
}

func (bs *BannersService) CreateBanner(banner entities.Banner) (int, error) {
	return bs.repo.CreateBanner(banner)
}

func (bs *BannersService) GetAllBanners(filters entities.QueryFilters) ([]entities.Banner, error) {
	return bs.repo.GetAllBanners(filters)
}

func (bs *BannersService) DeleteBanners(filters entities.QueryFilters) error {
	return bs.repo.DeleteBanners(filters)
}

func (bs *BannersService) UserBanner(filters entities.QueryFilters) (entities.Content, error) {
	return bs.repo.UserBanner(filters)
}
