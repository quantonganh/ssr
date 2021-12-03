package postgresql

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/quantonganh/ssr"
)

type repositoryService struct {
	db *gorm.DB
}

func NewRepositoryService(db *gorm.DB) ssr.RepositoryService {
	return &repositoryService{
		db: db,
	}
}

func (s *repositoryService) Create(r *ssr.Repository) error {
	if err := s.db.Create(&r).Error; err != nil {
		return errors.Wrapf(err, "failed to create repository: %s", r.FullName)
	}
	return nil
}

func (s *repositoryService) Get(id uint64) (*ssr.Repository, error) {
	var repo ssr.Repository
	s.db.First(&repo, "id = ?", id)
	return &repo, nil
}
