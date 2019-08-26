package http_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
	orderHttp "github.com/kodersky/golang-api-example/internal/app/api/order/delivery/http"
	"github.com/kodersky/golang-api-example/internal/app/api/order/mocks"
)

// TestFetch checks if http code is 200 and json response has correct format
func TestFetch(t *testing.T) {
	var mockOrder models.Order
	uid := uuid.New()
	mockOrder.ID = int64(10)
	mockOrder.UUID = uid
	mockOrder.Status = 0

	mockUCase := new(mocks.Usecase)
	mockListOrder := make([]*models.Order, 0)
	mockListOrder = append(mockListOrder, &mockOrder)
	limit := 10
	offset := 0
	mockUCase.On("Fetch", mock.Anything, limit, offset).Return(mockListOrder, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/orders", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.FetchOrder(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[{"distance":0, "id":"`+uid.String()+`", "status":"UNASSIGNED"}]`, rec.Body.String())
	//assert.
	mockUCase.AssertExpectations(t)
}

// TestFetchNoResults checks if response is empty array
func TestFetchNoResults(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	mockListOrder := make([]*models.Order, 0)
	limit := 1
	page := 100
	offset := limit * (page - 1)
	mockUCase.On("Fetch", mock.Anything, limit, offset).Return(mockListOrder, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/orders?page=100&limit=1", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.FetchOrder(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
	mockUCase.AssertExpectations(t)
}

// Test if http code is 500 and json response has correct format
func TestFetchError(t *testing.T) {
	mockUCase := new(mocks.Usecase)
	limit := 10
	offset := 0
	mockUCase.On("Fetch", mock.Anything, limit, offset).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/orders?page=1", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.FetchOrder(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.JSONEq(t, `{"error": "internal server error"}`, rec.Body.String())
	mockUCase.AssertExpectations(t)
}

// TestStore checks if response and json format are correct for newly added Order
func TestStore(t *testing.T) {
	//var googleClient *fakeClient
	mockOrder := struct {
		Origin      [2]string `json:"origin"`
		Destination [2]string `json:"destination"`
	}{
		Origin:      [2]string{"13.742310", "100.631418"},
		Destination: [2]string{"13.754179", "100.630377"},
	}

	tempMockOrder := mockOrder
	mockUCase := new(mocks.Usecase)

	j, err := json.Marshal(tempMockOrder)
	assert.NoError(t, err)

	mockUCase.On("Store", mock.Anything, mock.AnythingOfType("*models.Order")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/orders", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/orders")

	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"id":"00000000-0000-0000-0000-000000000000","status":"UNASSIGNED","distance":0}`, rec.Body.String())
	mockUCase.AssertExpectations(t)
}

// TestUpdate checks if order can be taken.
func TestUpdate(t *testing.T) {
	var mockOrder models.Order
	uid := uuid.New()
	mockOrder.ID = int64(10)
	mockOrder.UUID = uid
	mockOrder.Status = 0

	mockUCase := new(mocks.Usecase)

	mockUCase.On("GetByID", mock.Anything, uid.String()).Return(&mockOrder, nil)

	mockUCase.On("Update", mock.Anything, &mockOrder).Return(nil)

	e := echo.New()
	var jsonStr = []byte(`{"status":"TAKEN"}`)
	req, err := http.NewRequest(echo.PATCH, "/orders/"+uid.String(), bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("orders/:id")
	c.SetParamNames("id")
	c.SetParamValues(uid.String())
	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.Update(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)

}

// TestUpdateError checks if order cannot be taken with incorrect status.
func TestUpdateError(t *testing.T) {
	var mockOrder models.Order
	uid := uuid.New()
	mockOrder.ID = int64(10)
	mockOrder.UUID = uid
	mockOrder.Status = 0

	mockUCase := new(mocks.Usecase)

	e := echo.New()
	var jsonStr = []byte(`{"status":"I WILL NOT TAKE IT"}`)
	req, err := http.NewRequest(echo.PATCH, "/orders/"+uid.String(), bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("orders/:id")
	c.SetParamNames("id")
	c.SetParamValues(uid.String())
	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.Update(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	mockUCase.AssertExpectations(t)

}

// TestUpdateTakenError checks if order cannot be taken second time
func TestUpdateTakenError(t *testing.T) {
	var mockOrder models.Order
	uid := uuid.New()
	mockOrder.ID = int64(10)
	mockOrder.UUID = uid
	mockOrder.Status = 1

	mockUCase := new(mocks.Usecase)

	mockUCase.On("GetByID", mock.Anything, uid.String()).Return(&mockOrder, nil)

	mockUCase.On("Update", mock.Anything, &mockOrder).Return(models.ErrConflict)

	e := echo.New()
	var jsonStr = []byte(`{"status":"TAKEN"}`)
	req, err := http.NewRequest(echo.PATCH, "/orders/"+uid.String(), bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("orders/:id")
	c.SetParamNames("id")
	c.SetParamValues(uid.String())
	handler := orderHttp.OrderHandler{
		OUsecase: mockUCase,
	}
	err = handler.Update(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusConflict, rec.Code)
	mockUCase.AssertExpectations(t)
}
