// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	model "vod/model"

	mock "github.com/stretchr/testify/mock"
)

// IEncryption is an autogenerated mock type for the IEncryption type
type IEncryption struct {
	mock.Mock
}

// Encrypt provides a mock function with given fields: videoLinks
func (_m *IEncryption) Encrypt(videoLinks []model.VideoLinks) (map[model.VideoLinks]model.VideoLinks, error) {
	ret := _m.Called(videoLinks)

	var r0 map[model.VideoLinks]model.VideoLinks
	var r1 error
	if rf, ok := ret.Get(0).(func([]model.VideoLinks) (map[model.VideoLinks]model.VideoLinks, error)); ok {
		return rf(videoLinks)
	}
	if rf, ok := ret.Get(0).(func([]model.VideoLinks) map[model.VideoLinks]model.VideoLinks); ok {
		r0 = rf(videoLinks)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[model.VideoLinks]model.VideoLinks)
		}
	}

	if rf, ok := ret.Get(1).(func([]model.VideoLinks) error); ok {
		r1 = rf(videoLinks)
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
