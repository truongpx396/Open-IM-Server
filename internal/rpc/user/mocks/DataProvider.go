// Code generated by mockery v2.46.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DataProvider is an autogenerated mock type for the DataProvider type
type DataProvider struct {
	mock.Mock
}

// GetRandomNumber provides a mock function with given fields: id
func (_m *DataProvider) GetRandomNumber(id int) (int, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetRandomNumber")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (int, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) int); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDataProvider creates a new instance of DataProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataProvider {
	mock := &DataProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
