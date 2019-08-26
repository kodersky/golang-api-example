package usecase_test

import (
	//"github.com/kodersky/golang-api-example/internal/app/api/models"
	"context"
	"errors"
	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
	"github.com/kodersky/golang-api-example/internal/app/api/order/mocks"
	useCase "github.com/kodersky/golang-api-example/internal/app/api/order/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"googlemaps.github.io/maps"
	"testing"
	"time"
)

type fakeClient struct {
	gr maps.DistanceMatrixResponse
	t  time.Duration
	e  error
}

var googleResponse = maps.DistanceMatrixResponse{
	OriginAddresses:      []string{"1119 Rama IX Rd, Khwaeng Suan Luang, Khet Suan Luang, Krung Thep Maha Nakhon 10250, Thailand"},
	DestinationAddresses: []string{"2098/257 Soi Ramkhamhaeng 24 Yeak 30, Khwaeng Hua Mak, Khet Bang Kapi, Krung Thep Maha Nakhon 10240, Thailand"},
	Rows: []maps.DistanceMatrixElementsRow{
		{
			Elements: []*maps.DistanceMatrixElement{
				{
					Status:            "OK",
					Duration:          212000000000,
					DurationInTraffic: 0,
					Distance: maps.Distance{
						HumanReadable: "1.4 km",
						Meters:        1426,
					},
				},
			},
		},
	},
}

func (fk *fakeClient) DistanceMatrix(ctx context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
	if fk.t > 0 {
		time.Sleep(4 * time.Second)
	}
	if fk.e != nil {
		return &maps.DistanceMatrixResponse{}, errors.New("google api problem")
	}
	return &googleResponse, nil
}

func newWithClient(gr maps.DistanceMatrixResponse, t time.Duration, e error) *order.Client {
	var c fakeClient
	c.gr = gr
	c.t = t
	c.e = e
	return &order.Client{
		Client: &c,
	}
}

func TestStore(t *testing.T) {
	mockOrderRepo := new(mocks.Repository)
	var mockOrder models.Order

	t.Run("success", func(t *testing.T) {
		tempMockOrder := mockOrder
		mockOrderRepo.On("Store", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil).Once()

		gc := newWithClient(googleResponse, 0, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		assert.NoError(t, err)
		assert.Equal(t, mockOrder.ID, tempMockOrder.ID)
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("coordinates-not-found", func(t *testing.T) {
		tempMockOrder := mockOrder

		googleBadResponse := googleResponse

		googleBadResponse.Rows[0].Elements[0].Status = "not found"

		gc := newWithClient(googleBadResponse, 0, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		if assert.Error(t, err) {
			assert.Equal(t, models.ErrBadParamInput, err)
		}
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("response-timeout", func(t *testing.T) {
		tempMockOrder := mockOrder

		gc := newWithClient(googleResponse, 2*time.Second, nil)

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		if assert.Error(t, err) {
			assert.Equal(t, models.ErrTimeout, err)
		}
		mockOrderRepo.AssertExpectations(t)
	})
	t.Run("google-client-error", func(t *testing.T) {
		tempMockOrder := mockOrder

		gc := newWithClient(googleResponse, 0, errors.New(""))

		u := useCase.NewOrderUsecase(mockOrderRepo, 2*time.Second, gc)

		err := u.Store(context.TODO(), &tempMockOrder)

		assert.Equal(t, models.ErrInternalServerError, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
