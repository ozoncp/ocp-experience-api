package repo

import "github.com/ozoncp/ocp-experience-api/internal/models"

// Repo represents an experience storage
type Repo interface {
	Add(requests []models.Experience) error
	List(limit, offset uint64) ([]models.Experience, error)
	Describe(id uint64) (*models.Experience, error)
}
