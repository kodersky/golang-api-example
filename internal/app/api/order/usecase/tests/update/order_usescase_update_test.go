package order_usecase_update_test

import (
	"context"
	"errors"
	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order/mocks"
	ucase "github.com/kodersky/golang-api-example/internal/app/api/order/usecase"
	"github.com/kodersky/golang-api-example/internal/app/api/order/usecase/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	mockOrder := models.Order{
		Status: 0,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("Update", mock.Anything, &mockOrder).Once().Return(nil)

		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)

		err := u.Update(context.TODO(), &mockOrder)
		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("Update", mock.Anything, &mockOrder).Once().Return(errors.New("error"))

		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)

		err := u.Update(context.TODO(), &mockOrder)
		assert.Error(t, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
