// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ilya-rusyanov/shrinklator/internal/handlers (interfaces: DeleteService)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -package mocks -destination ./mocks/mock_delete_service.go . DeleteService
//
// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/ilya-rusyanov/shrinklator/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockDeleteService is a mock of DeleteService interface.
type MockDeleteService struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteServiceMockRecorder
}

// MockDeleteServiceMockRecorder is the mock recorder for MockDeleteService.
type MockDeleteServiceMockRecorder struct {
	mock *MockDeleteService
}

// NewMockDeleteService creates a new mock instance.
func NewMockDeleteService(ctrl *gomock.Controller) *MockDeleteService {
	mock := &MockDeleteService{ctrl: ctrl}
	mock.recorder = &MockDeleteServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteService) EXPECT() *MockDeleteServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleteService) Delete(arg0 context.Context, arg1 entities.DeleteRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleteServiceMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleteService)(nil).Delete), arg0, arg1)
}
