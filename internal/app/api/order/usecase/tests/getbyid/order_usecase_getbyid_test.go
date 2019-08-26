package usecase_getbyid__test

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

func TestGetByID(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	//var mockOrder *models.Order

	mockOrder := models.Order{
		Status: 0,
	}

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&mockOrder, nil).Once()
		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)

		o, err := u.GetByID(context.TODO(), mockOrder.UUID.String())

		assert.NoError(t, err)
		assert.NotNil(t, o)

		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("error")).Once()

		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)
		o, err := u.GetByID(context.TODO(), mockOrder.UUID.String())

		assert.Error(t, err)
		assert.Nil(t, o)

		mockOrderRepo.AssertExpectations(t)
	})

}
