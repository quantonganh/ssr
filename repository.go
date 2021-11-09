package ssr

type Repository struct {
	ID uint64 `json:"id"`
	Provider string `json:"provider"`
	FullName string `json:"full_name"`
	Description string `json:"description"`
}

type RepositoryService interface {
	Create(r *Repository) error
	Get(repoID uint64) (*Repository, error)
}