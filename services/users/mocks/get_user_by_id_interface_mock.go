// Code generated by MockGen. DO NOT EDIT.
// Source: services/users/app/get_user_by_id.go
//
// Generated by this command:
//
//	mockgen -source=services/users/app/get_user_by_id.go -destination=services/users/mocks/get_user_by_id_interface_mock.go -package=mocks -write_generate_directive
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	app "github.com/italoservio/braz_ecommerce/services/users/app"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -source=services/users/app/get_user_by_id.go -destination=services/users/mocks/get_user_by_id_interface_mock.go -package=mocks -write_generate_directive

// MockGetUserByIdInterface is a mock of GetUserByIdInterface interface.
type MockGetUserByIdInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGetUserByIdInterfaceMockRecorder
}

// MockGetUserByIdInterfaceMockRecorder is the mock recorder for MockGetUserByIdInterface.
type MockGetUserByIdInterfaceMockRecorder struct {
	mock *MockGetUserByIdInterface
}

// NewMockGetUserByIdInterface creates a new mock instance.
func NewMockGetUserByIdInterface(ctrl *gomock.Controller) *MockGetUserByIdInterface {
	mock := &MockGetUserByIdInterface{ctrl: ctrl}
	mock.recorder = &MockGetUserByIdInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetUserByIdInterface) EXPECT() *MockGetUserByIdInterfaceMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockGetUserByIdInterface) Do(id string) (*app.GetUserByIdOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", id)
	ret0, _ := ret[0].(*app.GetUserByIdOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockGetUserByIdInterfaceMockRecorder) Do(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockGetUserByIdInterface)(nil).Do), id)
}
