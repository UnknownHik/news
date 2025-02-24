package service

import (
	"news-rest-api/internal/dto"
	e "news-rest-api/internal/pkg/errors"
	"news-rest-api/internal/repository"
)

type Service interface {
	UpdateNews(news dto.News, id uint64) error
	GetNewsList(page, pageSize int) ([]dto.News, error)
}

type DefaultService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *DefaultService {
	return &DefaultService{repo: repo}
}

func (s *DefaultService) UpdateNews(news dto.News, id uint64) error {
	for _, cat := range news.Categories {
		if cat < 0 {
			return e.ErrInvalidCategory
		}
	}
	return s.repo.UpdateNews(news, id)
}

func (s *DefaultService) GetNewsList(page, pageSize int) ([]dto.News, error) {
	offset := (page - 1) * pageSize
	return s.repo.GetNewsList(pageSize, offset)
}
