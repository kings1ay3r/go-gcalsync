// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "gcalsync/gophers/dto"

	mock "github.com/stretchr/testify/mock"
)

// Core is an autogenerated mock type for the Core type
type Core struct {
	mock.Mock
}

// GetAuthCodeURL provides a mock function with given fields: ctx
func (_m *Core) GetAuthCodeURL(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAuthCodeURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMyCalendarEvents provides a mock function with given fields: ctx
func (_m *Core) GetMyCalendarEvents(ctx context.Context) ([]dto.Calendar, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetMyCalendarEvents")
	}

	var r0 []dto.Calendar
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]dto.Calendar, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []dto.Calendar); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.Calendar)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertCalendars provides a mock function with given fields: ctx, code
func (_m *Core) InsertCalendars(ctx context.Context, code string) error {
	ret := _m.Called(ctx, code)

	if len(ret) == 0 {
		panic("no return value specified for InsertCalendars")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ProcessWebhook provides a mock function with given fields: ctx, resourceID
func (_m *Core) ProcessWebhook(ctx context.Context, resourceID string) error {
	ret := _m.Called(ctx, resourceID)

	if len(ret) == 0 {
		panic("no return value specified for ProcessWebhook")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, resourceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RenewExpiringWatches provides a mock function with given fields: ctx
func (_m *Core) RenewExpiringWatches(ctx context.Context) {
	_m.Called(ctx)
}

// NewCore creates a new instance of Core. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Core {
	mock := &Core{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
