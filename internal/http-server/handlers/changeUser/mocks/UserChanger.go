// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserChanger is an autogenerated mock type for the UserChanger type
type UserChanger struct {
	mock.Mock
}

// ChangeUser provides a mock function with given fields: addSeg, delSeg, id
func (_m *UserChanger) ChangeUser(addSeg []string, delSeg []string, id int) (string, error) {
	ret := _m.Called(addSeg, delSeg, id)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func([]string, []string, int) (string, error)); ok {
		return rf(addSeg, delSeg, id)
	}
	if rf, ok := ret.Get(0).(func([]string, []string, int) string); ok {
		r0 = rf(addSeg, delSeg, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func([]string, []string, int) error); ok {
		r1 = rf(addSeg, delSeg, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewUserChanger interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserChanger creates a new instance of UserChanger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserChanger(t mockConstructorTestingTNewUserChanger) *UserChanger {
	mock := &UserChanger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
