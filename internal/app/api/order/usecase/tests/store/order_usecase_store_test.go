package usecase_store_test

import (
	"context"
	"errors"
	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order/mocks"
	useCase "github.com/kodersky/golang-api-example/internal/app/api/order/usecase"
	"github.com/kodersky/golang-api-example/internal/app/api/order/usecase/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	var mockOrder models.Order

	t.Run("success", func(t *testing.T) {
		tempMockOrder := mockOrder
		mockOrderRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil).Once()

		gc := tests.NewWithClient(tests.GoogleResponse, 0, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		assert.NoError(t, err)
		assert.Equal(t, mockOrder.ID, tempMockOrder.ID)
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("coordinates-not-found", func(t *testing.T) {
		tempMockOrder := mockOrder

		googleBadResponse := tests.GoogleResponse

		googleBadResponse.Rows[0].Elements[0].Status = "not found"

		gc := tests.NewWithClient(googleBadResponse, 0, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		if assert.Error(t, err) {
			assert.Equal(t, models.ErrBadParamInput, err)
		}
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("response-timeout", func(t *testing.T) {
		tempMockOrder := mockOrder

		gc := tests.NewWithClient(tests.GoogleResponse, 2*time.Second, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		if assert.Error(t, err) {
			assert.Equal(t, models.ErrTimeout, err)
		}
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("google-client-error", func(t *testing.T) {
		tempMockOrder := mockOrder

		gc := tests.NewWithClient(tests.GoogleResponse, 0, errors.New(""))

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		assert.Equal(t, models.ErrInternalServerError, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
