// Code generated by MockGen. DO NOT EDIT.
// Source: packages/logger/logger.go
//
// Generated by this command:
//
//	mockgen -source=packages/logger/logger.go -destination=services/users/mocks/logger_interface_mock.go -package=mocks -write_generate_directive
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	logger "github.com/italoservio/braz_ecommerce/packages/logger"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -source=packages/logger/logger.go -destination=services/users/mocks/logger_interface_mock.go -package=mocks -write_generate_directive

// MockLoggerInterface is a mock of LoggerInterface interface.
type MockLoggerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerInterfaceMockRecorder
}

// MockLoggerInterfaceMockRecorder is the mock recorder for MockLoggerInterface.
type MockLoggerInterfaceMockRecorder struct {
	mock *MockLoggerInterface
}

// NewMockLoggerInterface creates a new mock instance.
func NewMockLoggerInterface(ctrl *gomock.Controller) *MockLoggerInterface {
	mock := &MockLoggerInterface{ctrl: ctrl}
	mock.recorder = &MockLoggerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerInterface) EXPECT() *MockLoggerInterfaceMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockLoggerInterface) Error(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", msg)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerInterfaceMockRecorder) Error(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLoggerInterface)(nil).Error), msg)
}

// Info mocks base method.
func (m *MockLoggerInterface) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerInterfaceMockRecorder) Info(msg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLoggerInterface)(nil).Info), msg)
}

// WithCtx mocks base method.
func (m *MockLoggerInterface) WithCtx(ctx context.Context) *logger.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithCtx", ctx)
	ret0, _ := ret[0].(*logger.Logger)
	return ret0
}

// WithCtx indicates an expected call of WithCtx.
func (mr *MockLoggerInterfaceMockRecorder) WithCtx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithCtx", reflect.TypeOf((*MockLoggerInterface)(nil).WithCtx), ctx)
}
