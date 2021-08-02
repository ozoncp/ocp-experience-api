package flusher

import (
	"github.com/ozoncp/ocp-experience-api/internal/models"
	"github.com/ozoncp/ocp-experience-api/internal/repo"
	"github.com/ozoncp/ocp-experience-api/internal/utils"
)

// Flusher adds experience items to a storage
type Flusher interface {
	Flush(entities []models.Experience) ([]models.Experience, error)
}

// NewFlusher creates a new Flusher instance that writes experience to storage
func NewFlusher(chunkSize uint, requestRepo repo.Repo, ) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		requestRepo: requestRepo,
	}
}

type flusher struct {
	chunkSize   uint
	requestRepo repo.Repo
}

// Flush stores a slice of Experiences into the Repo. It makes experiences by bulks of a certain size.
// May returns some remains items with an error
func (f *flusher) Flush(experiences []models.Experience) ([]models.Experience, error) {
	remains := make([]models.Experience, 0, f.chunkSize)
	bulks, err := utils.SplitExperienceToBulks(experiences, int(f.chunkSize))

	if err != nil {
		return nil, err
	}

	for index, bulk := range bulks {
		if len(bulk) == int(f.chunkSize) {
			addErr := f.requestRepo.Add(bulk)

			if addErr != nil {
				remains = append(remains, experiences[index * int(f.chunkSize):]...)
				return remains, addErr
			}
		} else {
			remains = append(remains, bulk...) // last bulk should be kept in buffer
		}
	}

	return remains, nil
}
