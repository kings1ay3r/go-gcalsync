// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	oauth2 "golang.org/x/oauth2"
)

// TokenSource is an autogenerated mock type for the TokenSource type
type TokenSource struct {
	mock.Mock
}

// Token provides a mock function with given fields:
func (_m *TokenSource) Token() (*oauth2.Token, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Token")
	}

	var r0 *oauth2.Token
	var r1 error
	if rf, ok := ret.Get(0).(func() (*oauth2.Token, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *oauth2.Token); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTokenSource creates a new instance of TokenSource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTokenSource(t interface {
	mock.TestingT
	Cleanup(func())
}) *TokenSource {
	mock := &TokenSource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
