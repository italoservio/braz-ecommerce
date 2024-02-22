// Code generated by MockGen. DO NOT EDIT.
// Source: services/users/app/create_user.go
//
// Generated by this command:
//
//	mockgen -source=services/users/app/create_user.go -destination=services/users/mocks/create_user_interface_mock.go -package=mocks -write_generate_directive
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	app "github.com/italoservio/braz_ecommerce/services/users/app"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -source=services/users/app/create_user.go -destination=services/users/mocks/create_user_interface_mock.go -package=mocks -write_generate_directive

// MockCreateUserInterface is a mock of CreateUserInterface interface.
type MockCreateUserInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCreateUserInterfaceMockRecorder
}

// MockCreateUserInterfaceMockRecorder is the mock recorder for MockCreateUserInterface.
type MockCreateUserInterfaceMockRecorder struct {
	mock *MockCreateUserInterface
}

// NewMockCreateUserInterface creates a new mock instance.
func NewMockCreateUserInterface(ctrl *gomock.Controller) *MockCreateUserInterface {
	mock := &MockCreateUserInterface{ctrl: ctrl}
	mock.recorder = &MockCreateUserInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateUserInterface) EXPECT() *MockCreateUserInterfaceMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockCreateUserInterface) Do(ctx context.Context, input *app.CreateUserInput) (*app.CreateUserOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, input)
	ret0, _ := ret[0].(*app.CreateUserOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockCreateUserInterfaceMockRecorder) Do(ctx, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockCreateUserInterface)(nil).Do), ctx, input)
}
