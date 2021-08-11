package utils

import (
	"testing"
	"time"

	"github.com/ozoncp/ocp-experience-api/internal/models"
	"github.com/stretchr/testify/assert"
)

var splitExperienceToBulksData = []models.Experience{
	models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
	models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
	models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
	models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
	models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
	models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
}

// bulks split test wrapper
func splitExperienceToBulksTest(t *testing.T, batches int, expected [][]models.Experience) {
	res, err := SplitExperienceToBulks(splitExperienceToBulksData, batches)

	assert.Equal(t, err, nil)
	assert.Equal(t, res, expected)
}

//
// SplitExperienceToBulks
//
func TestSplitExperienceToBulks1(t *testing.T) {
	expected := [][]models.Experience{
		{models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1)},
		{models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2)},
		{models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3)},
		{models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4)},
		{models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5)},
		{models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6)},
	}

	splitExperienceToBulksTest(t, 1, expected)
}

func TestSplitExperienceToBulks2(t *testing.T) {
	expected := [][]models.Experience{
		{
			models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
			models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
		},
		{
			models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
			models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
		},
		{
			models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
			models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
		},
	}

	splitExperienceToBulksTest(t, 2, expected)
}

func TestSplitExperienceToBulks3(t *testing.T) {
	expected := [][]models.Experience{
		{
			models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
			models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
			models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
		},
		{
			models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
			models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
			models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
		},
	}

	splitExperienceToBulksTest(t, 3, expected)
}

func TestSplitExperienceToBulks4(t *testing.T) {
	expected := [][]models.Experience{
		{
			models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
			models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
			models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
			models.NewExperience(4, 4, 4, time.Time{}, time.Time{}, 4),
		},
		{
			models.NewExperience(5, 5, 5, time.Time{}, time.Time{}, 5),
			models.NewExperience(6, 6, 6, time.Time{}, time.Time{}, 6),
		},
	}

	splitExperienceToBulksTest(t, 4, expected)
}

// checks on error
func TestSplitExperienceToBulks5(t *testing.T) {
	_, err := BatchSplit(batchSplitData, len(batchSplitData)+1)
	assert.NotEqual(t, err, nil)
}

//
// ConvertExperienceToMap
//
func TestConvertExperienceToMap1(t *testing.T) {
	data := []models.Experience{
		models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
		models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
		models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
	}

	res, err := ConvertExperienceToMap(data)
	expected := map[uint64]models.Experience{
		1: models.NewExperience(1, 1, 1, time.Time{}, time.Time{}, 1),
		2: models.NewExperience(2, 2, 2, time.Time{}, time.Time{}, 2),
		3: models.NewExperience(3, 3, 3, time.Time{}, time.Time{}, 3),
	}

	assert.Equal(t, err, nil)
	assert.Equal(t, res, expected)
}

// checks on error
func TestConvertExperienceToMap2(t *testing.T) {
	var data []models.Experience = nil
	_, err := ConvertExperienceToMap(data)

	assert.NotEqual(t, err, nil)
}
