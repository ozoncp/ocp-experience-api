package utils

import (
	"errors"
	"math"

	"github.com/ozoncp/ocp-experience-api/internal/models"
)

// SplitExperienceToBulks splits entire Experience slice to same batches with batch size, except last batch,
func SplitExperienceToBulks(entities []models.Experience, batchSize int) ([][]models.Experience, error) {
	var entitiesSize = len(entities)

	if entitiesSize < batchSize {
		return nil, errors.New("entire Experience slice size is lower than batch size")
	}

	var batchesCount = int(math.Ceil(float64(entitiesSize) / float64(batchSize)))
	var result = make([][]models.Experience, 0, batchesCount)

	for i := 0; i < entitiesSize; i = i + batchSize {
		var end = i + batchSize

		if end > entitiesSize {
			end = entitiesSize
		}

		result = append(result, entities[i:end])
	}

	return result, nil
}

// ConvertExperienceToMap converts entire Experience slice to hash table [id, experience]
func ConvertExperienceToMap(entities []models.Experience) (map[uint64]models.Experience, error) {
	var entitiesSize = len(entities)

	if entitiesSize == 0 {
		return nil, errors.New("entire Experience slice is empty")
	}

	var res = make(map[uint64]models.Experience, len(entities))

	for i := 0; i < entitiesSize; i++ {
		var value = entities[i]
		res[value.Id] = value
	}

	return res, nil
}
