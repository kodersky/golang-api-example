package usecase_fetch_test

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

func TestFetch(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	//var mockOrder *models.Order

	mockOrder := &models.Order{
		Status: 0,
	}

	mockListOrder := make([]*models.Order, 0)
	mockListOrder = append(mockListOrder, mockOrder)

	t.Run("success", func(t *testing.T) {
		mockOrderRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(mockListOrder, nil).Once()

		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)
		limit := 10
		offset := 0
		list, err := u.Fetch(context.TODO(), limit, offset)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListOrder))

		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockOrderRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int"),
			mock.AnythingOfType("int")).Return(nil, errors.New("error")).Once()

		gm := tests.NewWithClient(tests.GoogleResponse, 0, nil)
		u := ucase.NewOrderUsecase(mockOrderRepo, time.Second*2, gm)
		limit := 1
		offset := 0
		list, err := u.Fetch(context.TODO(), limit, offset)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockOrderRepo.AssertExpectations(t)
	})

}
