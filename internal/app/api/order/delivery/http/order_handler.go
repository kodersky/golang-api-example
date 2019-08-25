package http

import (
	"context"
	"fmt"
	"github.com/kodersky/golang-api-example/internal/app/api/order/delivery/http/helpers"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
)

// OrderHandler represent the httphandler for order
type OrderHandler struct {
	OUsecase order.Usecase
}

// NewOrderHandler will initialize the orders/ resources endpoint
func NewOrderHandler(e *echo.Echo, ou order.Usecase) {
	handler := &OrderHandler{
		OUsecase: ou,
	}
	e.GET("/orders", handler.FetchOrder)
	e.POST("/orders", handler.Store)
	e.PATCH("/orders/:id", handler.Update)
}

// FetchOrder will fetch the orders based on given params
func (o *OrderHandler) FetchOrder(c echo.Context) error {

	var pagination helpers.Pagination

	// We could use also Validator.v9 for this
	err := helpers.IsPaginationValid(&pagination, c.QueryParam("page"), c.QueryParam("limit"))

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: err.Error()})
	}

	offset := pagination.Limit * (pagination.Page - 1)

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	listOr, err := o.OUsecase.Fetch(ctx, pagination.Limit, offset)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, listOr)
}

// Update will update Order status
func (o *OrderHandler) Update(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type status struct {
		Status string
	}

	var s status

	err := c.Bind(&s)
	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: models.ErrBadParamInput.Error()})
	}

	if s.Status != "TAKEN" {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: models.ErrBadParamInput.Error()})
	}

	or, err := o.OUsecase.GetByID(ctx, id)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: err.Error()})
	}

	err = o.OUsecase.Update(ctx, or)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: err.Error()})
	}

	logrus.Infof(fmt.Sprintf("Order %s status changed to %s", or.UUID, or.Status))

	return c.JSON(http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "SUCCESS"})
}

// Store will store the order by given request body
func (o *OrderHandler) Store(c echo.Context) error {
	var or models.Order
	var orderRequest helpers.OrderStruct

	err := c.Bind(&orderRequest)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, helpers.ResponseError{Message: models.ErrBadParamInput.Error()})
	}

	if ok, _ := helpers.IsOrderReqValid(&orderRequest); !ok {
		return c.JSON(http.StatusUnprocessableEntity, helpers.ResponseError{Message: models.ErrBadParamInput.Error()})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// err's can be skipped because we've validated lat, long earlier
	or.StartLat, _ = strconv.ParseFloat(orderRequest.Origin[0], 64)
	or.StartLong, _ = strconv.ParseFloat(orderRequest.Origin[1], 64)
	or.EndLat, _ = strconv.ParseFloat(orderRequest.Destination[0], 64)
	or.EndLong, _ = strconv.ParseFloat(orderRequest.Destination[1], 64)

	err = o.OUsecase.Store(ctx, &or)

	if err != nil {
		return c.JSON(helpers.GetStatusCode(err), helpers.ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, or)
}
