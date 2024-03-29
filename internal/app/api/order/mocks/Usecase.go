// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/kodersky/golang-api-example/internal/app/api/models"
	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx, limit, offset
func (_m *Usecase) Fetch(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	ret := _m.Called(ctx, limit, offset)

	var r0 []*models.Order
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []*models.Order); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Usecase) GetByID(ctx context.Context, id string) (*models.Order, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Order
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Order); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, or
func (_m *Usecase) Store(ctx context.Context, or *models.Order) error {
	ret := _m.Called(ctx, or)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Order) error); ok {
		r0 = rf(ctx, or)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, or
func (_m *Usecase) Update(ctx context.Context, or *models.Order) error {
	ret := _m.Called(ctx, or)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Order) error); ok {
		r0 = rf(ctx, or)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
