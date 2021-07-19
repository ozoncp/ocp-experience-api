package utils

import (
	"reflect"
	"testing"
)

var batchSplitData = []string{"hello", "world", "golang", "school", "test", "utils"}
var filterSliceData = batchSplitData	// change if you need

// checks slice or map on equality
func assertEqual(t *testing.T, current, expected interface{}) {
	if !reflect.DeepEqual(current, expected) {
		t.Errorf("Current and expected entities are not equal. Current %v, expected %v", current, expected)
	}
}

// batch split test wrapper
func batchSplitTest(t *testing.T, batches int, expected [][]string) {
	var res, err = BatchSplit(batchSplitData, batches)

	if err != nil {
		t.Errorf("Batch split error, %s", err.Error())
	}

	assertEqual(t, res, expected)
}

//
// BatchSplit
//
func TestBatchSplit1(t *testing.T) {
	var expected = [][]string{{"hello"}, {"world"}, {"golang"}, {"school"}, {"test"}, {"utils"}}
	batchSplitTest(t, 1, expected)
}

func TestBatchSplit2(t *testing.T) {
	var expected = [][]string{{"hello", "world"}, {"golang", "school"}, {"test", "utils"}}
	batchSplitTest(t, 2, expected)
}

func TestBatchSplit3(t *testing.T) {
	var expected = [][]string{{"hello", "world", "golang"}, {"school", "test", "utils"}}
	batchSplitTest(t, 3, expected)
}

func TestBatchSplit4(t *testing.T) {
	var expected = [][]string{{"hello", "world", "golang", "school"}, {"test", "utils"}}
	batchSplitTest(t, 4, expected)
}

//
// ReverseMap
//
func TestReverseMap(t *testing.T) {
	var data = map[string]string {"key1" : "value1"}
	var res = ReverseMap(data)
	var expected = map[string]string {"value1" : "key1"}

	assertEqual(t, res, expected)
}

//
// FilterSlice
//
func TestFilterSlice1(t *testing.T) {
	var res = FilterSlice(filterSliceData, []string{"hello"})
	var expected = []string{"world", "golang", "school", "test", "utils"}

	assertEqual(t, res, expected)
}

func TestFilterSlice2(t *testing.T) {
	var res = FilterSlice(filterSliceData, []string{"hello", "golang"})
	var expected = []string{"world", "school", "test", "utils"}

	assertEqual(t, res, expected)
}

func TestFilterSlice3(t *testing.T) {
	var res = FilterSlice(filterSliceData, []string{"hello", "golang", "utils"})
	var expected = []string{"world", "school", "test"}

	assertEqual(t, res, expected)
}
