package postgresql

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/quantonganh/ssr"
)

const (
	sqlInsertRepository = `INSERT INTO repository (provider, full_name, description) VALUES ($1, $2, $3)`
	sqlSelectRepository = `SELECT id, provider, full_name, description FROM repository WHERE id = $1`
)

type repositoryService struct {
	db *sql.DB
}

func NewRepositoryService(db *sql.DB) ssr.RepositoryService {
	return &repositoryService{
		db: db,
	}
}

func (s *repositoryService) Create(r *ssr.Repository) error {
	_, err := s.db.Exec(sqlInsertRepository, r.Provider, r.FullName, r.Description)
	if err != nil {
		return errors.Wrapf(err, "failed to create repository: %s", r.FullName)
	}
	return nil
}

func (s *repositoryService) Get(id uint64) (*ssr.Repository, error) {
	var (
		provider string
		fullName string
		description string
	)
	err := s.db.QueryRow(sqlSelectRepository, id).Scan(&provider, &fullName, &description)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get repository: %d", id)
	}
	return &ssr.Repository{
		ID:          id,
		Provider:    provider,
		FullName:    fullName,
		Description: description,
	}, nil
}
