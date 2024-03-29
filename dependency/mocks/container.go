// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	types "github.com/Drafteame/inject/types"
	mock "github.com/stretchr/testify/mock"
)

// Container is an autogenerated mock type for the Container type
type Container struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *Container) Get(name types.Symbol) (interface{}, error) {
	ret := _m.Called(name)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(types.Symbol) interface{}); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(types.Symbol) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewContainer interface {
	mock.TestingT
	Cleanup(func())
}

// NewContainer creates a new instance of Container. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewContainer(t mockConstructorTestingTNewContainer) *Container {
	mock := &Container{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
