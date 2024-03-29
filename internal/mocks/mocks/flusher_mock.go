// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozoncp/ocp-experience-api/internal/flusher (interfaces: Flusher)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/ozoncp/ocp-experience-api/internal/models"
)

// MockFlusher is a mock of Flusher interface.
type MockFlusher struct {
	ctrl     *gomock.Controller
	recorder *MockFlusherMockRecorder
}

// MockFlusherMockRecorder is the mock recorder for MockFlusher.
type MockFlusherMockRecorder struct {
	mock *MockFlusher
}

// NewMockFlusher creates a new mock instance.
func NewMockFlusher(ctrl *gomock.Controller) *MockFlusher {
	mock := &MockFlusher{ctrl: ctrl}
	mock.recorder = &MockFlusherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFlusher) EXPECT() *MockFlusherMockRecorder {
	return m.recorder
}

// Flush mocks base method.
func (m *MockFlusher) Flush(arg0 context.Context, arg1 []models.Experience) ([]models.Experience, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush", arg0, arg1)
	ret0, _ := ret[0].([]models.Experience)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Flush indicates an expected call of Flush.
func (mr *MockFlusherMockRecorder) Flush(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockFlusher)(nil).Flush), arg0, arg1)
}
