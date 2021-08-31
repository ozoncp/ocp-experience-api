package utils

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
)

var batchSplitData = []string{"hello", "world", "golang", "school", "test", "utils"}
var filterSliceData = batchSplitData // change if you need

// checks slice or map on equality
func assertEqual(t *testing.T, current, expected interface{}) {
	if !reflect.DeepEqual(current, expected) {
		t.Errorf("Current and expected entities are not equal. Current %v, expected %v", current, expected)
	}
}

// batch split test wrapper
func batchSplitTest(t *testing.T, batches int, expected [][]string) {
	res, err := BatchSplit(batchSplitData, batches)

	if err != nil {
		t.Errorf("Batch split error, %s", err.Error())
	}

	assertEqual(t, res, expected)
}

//
// BatchSplit
//
func TestBatchSplit1(t *testing.T) {
	expected := [][]string{{"hello"}, {"world"}, {"golang"}, {"school"}, {"test"}, {"utils"}}
	batchSplitTest(t, 1, expected)
}

func TestBatchSplit2(t *testing.T) {
	expected := [][]string{{"hello", "world"}, {"golang", "school"}, {"test", "utils"}}
	batchSplitTest(t, 2, expected)
}

func TestBatchSplit3(t *testing.T) {
	expected := [][]string{{"hello", "world", "golang"}, {"school", "test", "utils"}}
	batchSplitTest(t, 3, expected)
}

func TestBatchSplit4(t *testing.T) {
	expected := [][]string{{"hello", "world", "golang", "school"}, {"test", "utils"}}
	batchSplitTest(t, 4, expected)
}

// checks on error
func TestBatchSplit5(t *testing.T) {
	_, err := BatchSplit(batchSplitData, len(batchSplitData)+1)

	if err == nil {
		t.Errorf("Batch split err is not nil")
	}
}

//
// ReverseMap
//
func TestReverseMap(t *testing.T) {
	data := map[string]string{"key1": "value1"}
	res := ReverseMap(data)
	expected := map[string]string{"value1": "key1"}

	assertEqual(t, res, expected)
}

//
// FilterSlice
//
func TestFilterSlice1(t *testing.T) {
	res := FilterSlice(filterSliceData, []string{"hello"})
	expected := []string{"world", "golang", "school", "test", "utils"}

	assertEqual(t, res, expected)
}

func TestFilterSlice2(t *testing.T) {
	res := FilterSlice(filterSliceData, []string{"hello", "golang"})
	expected := []string{"world", "school", "test", "utils"}

	assertEqual(t, res, expected)
}

func TestFilterSlice3(t *testing.T) {
	res := FilterSlice(filterSliceData, []string{"hello", "golang", "utils"})
	expected := []string{"world", "school", "test"}

	assertEqual(t, res, expected)
}

// checks on error if data is empty
func TestFilterSlice4(t *testing.T) {
	res := FilterSlice(nil, []string{"hello", "golang", "utils"})

	if res != nil {
		t.Errorf("FilterSlice res should be nil if filter data is empty")
	}
}

// checks on the same data if filter is empty
func TestFilterSlice5(t *testing.T) {
	res := FilterSlice(filterSliceData, nil)

	assertEqual(t, res, res)
}

//
// ReadFileInLoop
//
const loopFileName = "./test_loop_file.txt"

// creates and writes test file by file name (file path)
func createTestLoopFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = file.WriteString("TestInformation")

	if err != nil {
		return err
	}

	return nil
}

// main TestReadFileInLoop body
func testReadFileInLoop(t *testing.T, fileName string, count int) {
	err := createTestLoopFile(fileName)

	if err != nil {
		log.Fatalf("Can not create file, err %v", err.Error())
	}

	readErr := ReadFileInLoop(fileName, count)
	assert.Equal(t, readErr, nil)

	removeErr := os.Remove(fileName)

	if removeErr != nil {
		log.Fatalf("Can not remove file, %v", err.Error())
	}
}

func TestReadFileInLoop1(t *testing.T) {
	testReadFileInLoop(t, loopFileName, 1)
}

func TestReadFileInLoop2(t *testing.T) {
	testReadFileInLoop(t, loopFileName, 5)
}
