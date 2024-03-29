package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

// BatchSplit splits entire slice to same batches with batch size, except last batch,
// returns error if can not split
func BatchSplit(slice []string, batchSize int) ([][]string, error) {
	size := len(slice)

	if size < batchSize {
		return nil, errors.New("entire slice size is lower than batch size")
	}

	batchesCount := int(math.Ceil(float64(len(slice)) / float64(batchSize)))
	result := make([][]string, 0, batchesCount)

	for i := 0; i < len(slice); i = i + batchSize {
		end := i + batchSize

		if end > len(slice) {
			end = len(slice)
		}

		result = append(result, slice[i:end])
	}

	return result, nil
}

// ReverseMap Swaps key and value at m map
func ReverseMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))

	for key, value := range m {
		result[value] = key
	}

	return result
}

// FilterSlice filters entire slice with filter
func FilterSlice(in, filter []string) []string {
	if len(in) == 0 {
		return nil
	}

	if len(filter) == 0 {
		return in
	}

	containsFunc := func(filter []string, value string) bool {
		for i := 0; i < len(filter); i++ {
			if filter[i] == value {
				return true
			}
		}

		return false
	}

	result := make([]string, 0, cap(in))

	for i := 0; i < len(in); i++ {
		value := in[i]

		if !containsFunc(filter, value) {
			result = append(result, value)
		}
	}

	return result
}

// ReadFileInLoop reads file in loop count arg times
func ReadFileInLoop(filePath string, count int) error {
	readFileFunc := func() ([]byte, error) {
		file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModeType)

		if err != nil {
			return nil, err
		}

		defer func() {
			err := file.Close()

			if err != nil {
				log.Fatal(err)
			}
		}()

		data, readErr := io.ReadAll(file)

		if readErr != nil {
			return nil, err
		}

		return data, nil
	}

	for i := 0; i < count; i++ {
		data, err := readFileFunc()

		if err != nil {
			return err
		}

		fmt.Printf("ReadFileInLoop iteration: %v, data: %v\n", i, string(data))
	}

	return nil
}
