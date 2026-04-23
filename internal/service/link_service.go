package service

import (
	"errors"
	"link-storage-service/internal/cache"
	"link-storage-service/internal/model"
	"link-storage-service/internal/repository"
)

type LinkService interface {
	Create(url string) (string, error)
	Get(code string) (*model.Link, error)
	GetAll(limit, offset int) ([]model.Link, error)
	Delete(code string) error
	GetStats(code string) (*model.Link, error)
}

type linkService struct {
	repo  repository.LinkRepository
	cache cache.Cache
}

func NewLinkService(repo repository.LinkRepository, cache cache.Cache) LinkService {
	return &linkService{
		repo:  repo,
		cache: cache,
	}
}

func (s *linkService) Create(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}

	link := &model.Link{
		OriginalURL: url,
	}

	err := s.repo.Create(link)
	if err != nil {
		return "", err
	}

	shortCode := encodeBase62(link.ID)

	err = s.repo.UpdateShortCode(link.ID, shortCode)
	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *linkService) Get(code string) (*model.Link, error) {
	visits, err := s.repo.IncrementVisits(code)
	if err != nil {
		return nil, err
	}

	if url, err := s.cache.Get(code); err == nil && url != "" {
		return &model.Link{
			ShortCode:   code,
			OriginalURL: url,
			Visits:      visits,
		}, nil
	}

	link, err := s.repo.GetByShortCode(code)
	if err != nil {
		return nil, err
	}

	link.Visits = visits

	_ = s.cache.Set(code, link.OriginalURL)

	return link, nil
}

func (s *linkService) GetAll(limit, offset int) ([]model.Link, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *linkService) Delete(code string) error {
	err := s.repo.Delete(code)
	if err != nil {
		return err
	}

	s.cache.Delete(code)

	return nil
}

func (s *linkService) GetStats(code string) (*model.Link, error) {
	return s.repo.GetByShortCode(code)
}
