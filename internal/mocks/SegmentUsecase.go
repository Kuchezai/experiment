// Code generated by mockery v2.33.0. DO NOT EDIT.

package mocks

import (
	entity "experiment.io/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// SegmentUsecase is an autogenerated mock type for the SegmentUsecase type
type SegmentUsecase struct {
	mock.Mock
}

// DeleteSegment provides a mock function with given fields: slug
func (_m *SegmentUsecase) DeleteSegment(slug string) error {
	ret := _m.Called(slug)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(slug)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSegment provides a mock function with given fields: seg
func (_m *SegmentUsecase) NewSegment(seg entity.Segment) error {
	ret := _m.Called(seg)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.Segment) error); ok {
		r0 = rf(seg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewSegmentWithAutoAssign provides a mock function with given fields: seg, percentAssigned
func (_m *SegmentUsecase) NewSegmentWithAutoAssign(seg entity.Segment, percentAssigned int) ([]int, error) {
	ret := _m.Called(seg, percentAssigned)

	var r0 []int
	var r1 error
	if rf, ok := ret.Get(0).(func(entity.Segment, int) ([]int, error)); ok {
		return rf(seg, percentAssigned)
	}
	if rf, ok := ret.Get(0).(func(entity.Segment, int) []int); ok {
		r0 = rf(seg, percentAssigned)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	if rf, ok := ret.Get(1).(func(entity.Segment, int) error); ok {
		r1 = rf(seg, percentAssigned)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSegmentUsecase creates a new instance of SegmentUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSegmentUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *SegmentUsecase {
	mock := &SegmentUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
