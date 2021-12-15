package ssr

type Repository struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
	Provider string `json:"provider"`
	FullName string `json:"full_name"`
	Description string `json:"description"`
}

func (Repository) TableName() string {
	return "repository"
}

type RepositoryService interface {
	Create(r *Repository) error
	Get(repoID uint64) (*Repository, error)
}