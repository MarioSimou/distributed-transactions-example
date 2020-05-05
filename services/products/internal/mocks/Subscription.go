// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	rabbitmq "products/internal/rabbitmq"

	mock "github.com/stretchr/testify/mock"
)

// Subscription is an autogenerated mock type for the Subscription type
type Subscription struct {
	mock.Mock
}

// GetHandlerFunc provides a mock function with given fields:
func (_m *Subscription) GetHandlerFunc() rabbitmq.HandlerFunc {
	ret := _m.Called()

	var r0 rabbitmq.HandlerFunc
	if rf, ok := ret.Get(0).(func() rabbitmq.HandlerFunc); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rabbitmq.HandlerFunc)
		}
	}

	return r0
}

// GetQueue provides a mock function with given fields:
func (_m *Subscription) GetQueue() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
