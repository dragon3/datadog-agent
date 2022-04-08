// Code generated by mockery v2.10.2. DO NOT EDIT.

package mocks

import (
	eval "github.com/DataDog/datadog-agent/pkg/compliance/eval"
	mock "github.com/stretchr/testify/mock"
)

// Configuration is an autogenerated mock type for the Configuration type
type Configuration struct {
	mock.Mock
}

// EtcGroupPath provides a mock function with given fields:
func (_m *Configuration) EtcGroupPath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EvaluateFromCache provides a mock function with given fields: e
func (_m *Configuration) EvaluateFromCache(e eval.Evaluatable) (interface{}, error) {
	ret := _m.Called(e)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(eval.Evaluatable) interface{}); ok {
		r0 = rf(e)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(eval.Evaluatable) error); ok {
		r1 = rf(e)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hostname provides a mock function with given fields:
func (_m *Configuration) Hostname() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsLeader provides a mock function with given fields:
func (_m *Configuration) IsLeader() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MaxEventsPerRun provides a mock function with given fields:
func (_m *Configuration) MaxEventsPerRun() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// NodeLabels provides a mock function with given fields:
func (_m *Configuration) NodeLabels() map[string]string {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// NormalizeToHostRoot provides a mock function with given fields: path
func (_m *Configuration) NormalizeToHostRoot(path string) string {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RelativeToHostRoot provides a mock function with given fields: path
func (_m *Configuration) RelativeToHostRoot(path string) string {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
