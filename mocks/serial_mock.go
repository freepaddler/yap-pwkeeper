// Code generated by MockGen. DO NOT EDIT.
// Source: serial.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// GetSerials mocks base method.
func (m *MockSource) GetSerials(ctx context.Context, n int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSerials", ctx, n)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSerials indicates an expected call of GetSerials.
func (mr *MockSourceMockRecorder) GetSerials(ctx, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSerials", reflect.TypeOf((*MockSource)(nil).GetSerials), ctx, n)
}