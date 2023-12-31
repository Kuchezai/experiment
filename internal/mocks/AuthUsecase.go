// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	entity "experiment.io/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// AuthUsecase is an autogenerated mock type for the AuthUsecase type
type AuthUsecase struct {
	mock.Mock
}

// Login provides a mock function with given fields: user
func (_m *AuthUsecase) Login(user entity.User) (string, error) {
	ret := _m.Called(user)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.User) (string, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(entity.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(entity.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Registration provides a mock function with given fields: user
func (_m *AuthUsecase) Registration(user entity.User) (int, error) {
	ret := _m.Called(user)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.User) (int, error)); ok {
		return rf(user)
	}
	if rf, ok := ret.Get(0).(func(entity.User) int); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(entity.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthUsecase creates a new instance of AuthUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthUsecase {
	mock := &AuthUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
