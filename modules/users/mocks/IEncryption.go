// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IEncryption is an autogenerated mock type for the IEncryption type
type IEncryption struct {
	mock.Mock
}

// CompareHashAndPassword provides a mock function with given fields: passwordForLogin, passwordInDB
func (_m *IEncryption) CompareHashAndPassword(passwordForLogin string, passwordInDB string) error {
	ret := _m.Called(passwordForLogin, passwordInDB)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(passwordForLogin, passwordInDB)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GenerateFromPassword provides a mock function with given fields: password
func (_m *IEncryption) GenerateFromPassword(password string) (string, error) {
	ret := _m.Called(password)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(password)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIEncryption interface {
	mock.TestingT
	Cleanup(func())
}

// NewIEncryption creates a new instance of IEncryption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIEncryption(t mockConstructorTestingTNewIEncryption) *IEncryption {
	mock := &IEncryption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
