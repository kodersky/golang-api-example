package helpers

import (
	"github.com/kodersky/golang-api-example/internal/app/api/models"
	h "github.com/kodersky/golang-api-example/internal/app/api/order/delivery/http"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

func IsOrderReqValid(m *h.OrderStruct) (bool, error) {
	validate := validator.New()
	err := validate.RegisterValidation("geo", validateGeo)
	if err != nil {
		return false, err
	}

	err = validate.Struct(m)
	if err != nil {
		return false, models.ErrBadParamInput
	}
	return true, nil
}

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	case models.ErrTimeout:
		return http.StatusRequestTimeout
	case models.ErrBadParamInput:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

// validateGeo implements validator.Func
func validateGeo(fl validator.FieldLevel) bool {
	co, ok := fl.Field().Interface().([2]string)

	type geoCoordinates struct {
		Lat string `validate:"latitude"`
		Lng string `validate:"longitude"`
	}

	if !ok {
		return false
	}

	gc := geoCoordinates{
		Lat: co[0],
		Lng: co[1],
	}

	validate := validator.New()

	err := validate.Struct(gc)
	if err != nil {
		return false
	}

	return true
}

func IsPaginationValid(pagination *h.Pagination, pageS string, limitS string) error {
	if pageS != "" {
		p, err := strconv.Atoi(pageS)
		if err != nil || p <= 0 {
			return models.ErrBadParamInput
		}
		pagination.Page = p
	} else {
		pagination.Page = 1
	}

	if limitS != "" {
		l, err := strconv.Atoi(limitS)
		if err != nil || l < 0 {
			return models.ErrBadParamInput
		}
		pagination.Limit = l
	} else {
		pagination.Limit = 10
	}

	return nil
}
