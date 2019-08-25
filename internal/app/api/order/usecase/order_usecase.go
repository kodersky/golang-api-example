package usecase

import (
	"context"
	"fmt"
	"time"

	"googlemaps.github.io/maps"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
	"github.com/spf13/viper"
)

type orderUsecase struct {
	orderRepo      order.Repository
	contextTimeout time.Duration
}

// NewOrderUsecase will create new an orderUsecase object representation of order.Usecase interface
func NewOrderUsecase(o order.Repository, timeout time.Duration) order.Usecase {
	return &orderUsecase{
		orderRepo:      o,
		contextTimeout: timeout,
	}
}

func (o *orderUsecase) Fetch(c context.Context, limit int, offset int) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	listOrder, err := o.orderRepo.Fetch(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return listOrder, nil
}

func (o *orderUsecase) GetByID(c context.Context, id string) (*models.Order, error) {

	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	res, err := o.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (o *orderUsecase) Update(c context.Context, or *models.Order) error {

	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	switch os := or.Status; os {
	case models.Unassigned:
		or.Status = models.Taken
	case models.Taken:
		return models.ErrConflict
	default:
		return models.ErrInternalServerError
	}

	err := o.orderRepo.Update(ctx, or)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderUsecase) Store(c context.Context, m *models.Order) error {
	apiKey := viper.GetString(`gmaps_apikey`)

	ctx, cancel := context.WithTimeout(c, o.contextTimeout)
	defer cancel()

	c1 := make(chan *maps.DistanceMatrixResponse, 1)

	gm, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return models.ErrInternalServerError
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%.8f,%.8f", m.StartLat, m.StartLong)},
		Destinations: []string{fmt.Sprintf("%.8f,%.8f", m.EndLat, m.EndLong)},
	}

	go func() {
		// TODO: This err should be handled ;)
		route, _ := gm.DistanceMatrix(context.Background(), r)

		c1 <- route
	}()

	select {
	case res := <-c1:
		if res == nil {
			// Probably bad api key
			return models.ErrInternalServerError
		}
		if len(res.Rows) == 0 || res.Rows[0].Elements[0].Status != "OK" {
			return models.ErrBadParamInput
		}
		m.Distance = res.Rows[0].Elements[0].Distance.Meters
	case <-time.After(3 * time.Second):
		return models.ErrTimeout
	}

	err = o.orderRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
