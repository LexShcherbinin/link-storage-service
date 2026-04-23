package repository

import "link-storage-service/internal/model"

type LinkRepository interface {
	Create(link *model.Link) error
	GetByShortCode(code string) (*model.Link, error)
	GetAll(limit, offset int) ([]model.Link, error)
	Delete(code string) error
	IncrementVisits(code string) (int64, error)
	UpdateShortCode(id int64, code string) error
}
