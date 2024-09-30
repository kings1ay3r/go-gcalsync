// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"
	dao "gcalsync/gophers/dao"

	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"

	time "time"
)

// DAO is an autogenerated mock type for the DAO type
type DAO struct {
	mock.Mock
}

// FindCalendarByCalendarID provides a mock function with given fields: _a0, _a1
func (_m *DAO) FindCalendarByCalendarID(_a0 context.Context, _a1 string) (*dao.Calendar, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for FindCalendarByCalendarID")
	}

	var r0 *dao.Calendar
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*dao.Calendar, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *dao.Calendar); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.Calendar)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindExpiringWatches provides a mock function with given fields: _a0
func (_m *DAO) FindExpiringWatches(_a0 context.Context) ([]dao.WatchesWithDetails, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for FindExpiringWatches")
	}

	var r0 []dao.WatchesWithDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]dao.WatchesWithDetails, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []dao.WatchesWithDetails); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dao.WatchesWithDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserCalendars provides a mock function with given fields: _a0, _a1
func (_m *DAO) GetUserCalendars(_a0 context.Context, _a1 int) ([]dao.Calendar, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetUserCalendars")
	}

	var r0 []dao.Calendar
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]dao.Calendar, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []dao.Calendar); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dao.Calendar)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserTokens provides a mock function with given fields: _a0, _a1, _a2
func (_m *DAO) GetUserTokens(_a0 context.Context, _a1 int, _a2 string) (*oauth2.Token, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetUserTokens")
	}

	var r0 *oauth2.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) (*oauth2.Token, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *oauth2.Token); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertCalendar provides a mock function with given fields: _a0, _a1
func (_m *DAO) InsertCalendar(_a0 context.Context, _a1 dao.Calendar) (*dao.Calendar, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InsertCalendar")
	}

	var r0 *dao.Calendar
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dao.Calendar) (*dao.Calendar, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dao.Calendar) *dao.Calendar); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dao.Calendar)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dao.Calendar) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveUserCalendarData provides a mock function with given fields: _a0, _a1, _a2
func (_m *DAO) SaveUserCalendarData(_a0 context.Context, _a1 int, _a2 []*dao.CalendarData) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for SaveUserCalendarData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, []*dao.CalendarData) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUserTokens provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *DAO) SaveUserTokens(_a0 context.Context, _a1 int, _a2 string, _a3 string, _a4 string, _a5 time.Time) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	if len(ret) == 0 {
		panic("no return value specified for SaveUserTokens")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string, string, string, time.Time) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveWatch provides a mock function with given fields: _a0, _a1
func (_m *DAO) SaveWatch(_a0 context.Context, _a1 *dao.Watch) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SaveWatch")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *dao.Watch) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WatchExists provides a mock function with given fields: _a0, _a1, _a2
func (_m *DAO) WatchExists(_a0 context.Context, _a1 uint, _a2 string) (bool, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for WatchExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) (bool, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) bool); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDAO creates a new instance of DAO. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDAO(t interface {
	mock.TestingT
	Cleanup(func())
}) *DAO {
	mock := &DAO{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
