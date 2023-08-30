// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	entity "experiment.io/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// AddUserSegments provides a mock function with given fields: userID, added
func (_m *UserUsecase) AddUserSegments(userID int, added []entity.SlugWithExpiredDate) error {
	ret := _m.Called(userID, added)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []entity.SlugWithExpiredDate) error); ok {
		r0 = rf(userID, added)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveUserSegments provides a mock function with given fields: userID, removed
func (_m *UserUsecase) RemoveUserSegments(userID int, removed []string) error {
	ret := _m.Called(userID, removed)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, []string) error); ok {
		r0 = rf(userID, removed)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserSegments provides a mock function with given fields: userID
func (_m *UserUsecase) UserSegments(userID int) ([]entity.SlugWithExpiredDate, error) {
	ret := _m.Called(userID)

	var r0 []entity.SlugWithExpiredDate
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]entity.SlugWithExpiredDate, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(int) []entity.SlugWithExpiredDate); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.SlugWithExpiredDate)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UsersHistoryInCSVByDate provides a mock function with given fields: year, month
func (_m *UserUsecase) UsersHistoryInCSVByDate(year int, month int) (string, error) {
	ret := _m.Called(year, month)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) (string, error)); ok {
		return rf(year, month)
	}
	if rf, ok := ret.Get(0).(func(int, int) string); ok {
		r0 = rf(year, month)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(year, month)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserUsecase creates a new instance of UserUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserUsecase {
	mock := &UserUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
