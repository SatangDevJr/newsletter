// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	entity "newsletter/src/pkg/entity"
	error "newsletter/src/pkg/utils/error"

	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// GetAllSubscribers provides a mock function with given fields:
func (_m *UseCase) GetAllSubscribers() ([]entity.Subscribers, *error.ErrorCode) {
	ret := _m.Called()

	var r0 []entity.Subscribers
	var r1 *error.ErrorCode
	if rf, ok := ret.Get(0).(func() ([]entity.Subscribers, *error.ErrorCode)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []entity.Subscribers); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Subscribers)
		}
	}

	if rf, ok := ret.Get(1).(func() *error.ErrorCode); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*error.ErrorCode)
		}
	}

	return r0, r1
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
