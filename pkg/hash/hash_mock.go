// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/hash/hash.go

// Package hash is a generated GoMock package.
package hash

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHashMethod is a mock of HashMethod interface.
type MockHashMethod struct {
	ctrl     *gomock.Controller
	recorder *MockHashMethodMockRecorder
}

// MockHashMethodMockRecorder is the mock recorder for MockHashMethod.
type MockHashMethodMockRecorder struct {
	mock *MockHashMethod
}

// NewMockHashMethod creates a new mock instance.
func NewMockHashMethod(ctrl *gomock.Controller) *MockHashMethod {
	mock := &MockHashMethod{ctrl: ctrl}
	mock.recorder = &MockHashMethodMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHashMethod) EXPECT() *MockHashMethodMockRecorder {
	return m.recorder
}

// CompareValue mocks base method.
func (m *MockHashMethod) CompareValue(arg0, arg1 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareValue", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CompareValue indicates an expected call of CompareValue.
func (mr *MockHashMethodMockRecorder) CompareValue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareValue", reflect.TypeOf((*MockHashMethod)(nil).CompareValue), arg0, arg1)
}

// HashValue mocks base method.
func (m *MockHashMethod) HashValue(arg0 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashValue", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HashValue indicates an expected call of HashValue.
func (mr *MockHashMethodMockRecorder) HashValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashValue", reflect.TypeOf((*MockHashMethod)(nil).HashValue), arg0)
}
