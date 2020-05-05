// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"
)

// ChannelInterface is an autogenerated mock type for the ChannelInterface type
type ChannelInterface struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *ChannelInterface) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Consume provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *ChannelInterface) Consume(_a0 string, _a1 string, _a2 bool, _a3 bool, _a4 bool, _a5 bool, _a6 amqp.Table) (<-chan amqp.Delivery, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 <-chan amqp.Delivery
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, bool, bool, amqp.Table) <-chan amqp.Delivery); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan amqp.Delivery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, bool, bool, bool, bool, amqp.Table) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExchangeDeclare provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5, _a6
func (_m *ChannelInterface) ExchangeDeclare(_a0 string, _a1 string, _a2 bool, _a3 bool, _a4 bool, _a5 bool, _a6 amqp.Table) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5, _a6)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, bool, bool, amqp.Table) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5, _a6)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetQueueName provides a mock function with given fields:
func (_m *ChannelInterface) GetQueueName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Publish provides a mock function with given fields: exchange, key, mandatory, immediate, msg
func (_m *ChannelInterface) Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error {
	ret := _m.Called(exchange, key, mandatory, immediate, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, bool, bool, amqp.Publishing) error); ok {
		r0 = rf(exchange, key, mandatory, immediate, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueueBind provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *ChannelInterface) QueueBind(_a0 string, _a1 string, _a2 string, _a3 bool, _a4 amqp.Table) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, bool, amqp.Table) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueueDeclare provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4, _a5
func (_m *ChannelInterface) QueueDeclare(_a0 string, _a1 bool, _a2 bool, _a3 bool, _a4 bool, _a5 amqp.Table) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4, _a5)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, bool, bool, bool, bool, amqp.Table) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4, _a5)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
