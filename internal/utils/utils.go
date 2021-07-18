package utils

import "math"
import "errors"

// BatchSplit splits entire slice to same batches with batch size, except last batch,
// returns error if can not split
func BatchSplit(slice []string, batchSize int) ([][]string, error) {
	var size = len(slice)

	if size < batchSize {
		return nil, errors.New("entire slice size is lower than batch size")
	}

	var batchesCount = int(math.Ceil(float64(len(slice)) / float64(batchSize)))
	var result = make([][]string, 0, batchesCount)

	for i := 0; i < len(slice); i = i + batchSize {
		var end = i + batchSize

		if end > len(slice) {
			end = len(slice)
		}

		result = append(result, slice[i:end])
	}

	return result, nil
}

// ReverseMap Swaps key and value at m map
func ReverseMap(m map[string]string) map[string]string {
	var result = make(map[string]string, len(m))

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

	var containsFunc = func(filter []string, value string, ) bool {
		for i := 0; i < len(filter); i++ {
			if filter[i] == value {
				return true
			}
		}

		return false
	}

	var result = make([]string, 0, cap(in))

	for i := 0; i < len(in); i++ {
		var value = in[i]

		if !containsFunc(filter, value) {
			result = append(result, value)
		}
	}

	return result
}