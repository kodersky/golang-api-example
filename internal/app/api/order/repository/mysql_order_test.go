package repository_test

import (
	"context"
	"database/sql/driver"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
	orderRepo "github.com/kodersky/golang-api-example/internal/app/api/order/repository"
)

type any struct{}

func (a any) Match(v driver.Value) bool {
	return true
}

func TestFetch(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockOrders := []models.Order{
		{
			ID: 1, Status: 0, UUID: uuid.New(), Distance: 777,
			UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		{
			ID: 2, Status: 0, UUID: uuid.New(), Distance: 0,
			UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	columns := []string{"id", "uuid", "distance", "status", "updated_at", "created_at"}

	query := `SELECT id, uuid, distance, status, updated_at, created_at
	FROM orders
	INNER JOIN \(
	SELECT id
	FROM orders
	ORDER BY id
	LIMIT \? OFFSET \?\)
	AS my_orders USING\(id\)`

	mock.ExpectQuery(query).
		WithArgs(2, 0).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(mockOrders[0].ID, mockOrders[0].UUID, mockOrders[0].Distance,
			mockOrders[0].Status, mockOrders[0].UpdatedAt, mockOrders[0].CreatedAt).AddRow(mockOrders[1].ID, mockOrders[1].UUID, mockOrders[1].Distance,
			mockOrders[1].Status, mockOrders[1].UpdatedAt, mockOrders[1].CreatedAt))

	o := orderRepo.NewMysqlOrderRepository(db)
	// now we execute our method
	list, err := o.Fetch(context.TODO(), 2, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 2)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	uid := uuid.New()

	mockOrder := models.Order{
		ID: 1, Status: 0, UUID: uid, Distance: 777,
		UpdatedAt: time.Now(), CreatedAt: time.Now(),
	}

	columns := []string{"id", "uuid", "distance", "status", "updated_at", "created_at"}

	query := `SELECT id, uuid, distance, status, updated_at, created_at
 						FROM orders WHERE UUID = \?`
	assert.NoError(t, err)

	mock.ExpectQuery(query).
		WithArgs(any{}).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(mockOrder.ID, mockOrder.UUID, mockOrder.Distance,
			mockOrder.Status, mockOrder.UpdatedAt, mockOrder.CreatedAt))

	o := orderRepo.NewMysqlOrderRepository(db)
	// now we execute our method
	or, err := o.GetByID(context.TODO(), uid.String())

	assert.NoError(t, err)

	assert.NotNil(t, or)
	assert.Equal(t, or.ID, mockOrder.ID)
	assert.Equal(t, or.UUID, mockOrder.UUID)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStore(t *testing.T) {
	or := &models.Order{
		UUID:      uuid.New(),
		Status:    0,
		Distance:  1200,
		StartLat:  1,
		StartLong: 1,
		EndLat:    1,
		EndLong:   1,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := `INSERT orders SET distance=\?, uuid=\?, status=\?, start_lat=\?,
	start_long=\?, end_lat=\?, end_long=\?, updated_at=\?, created_at=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(or.Distance, any{}, or.Status, or.StartLat, or.StartLong, or.EndLat, or.EndLong, any{}, any{}).
		WillReturnResult(sqlmock.NewResult(12, 1))

	o := orderRepo.NewMysqlOrderRepository(db)

	err = o.Store(context.TODO(), or)
	assert.NoError(t, err)
	assert.NotNil(t, or)
	assert.Equal(t, int64(12), or.ID)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	or := &models.Order{
		ID:        int64(12),
		UUID:      uuid.New(),
		Status:    1,
		Distance:  1200,
		StartLat:  1,
		StartLong: 1,
		EndLat:    1,
		EndLong:   1,
		UpdatedAt: now,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE orders set status=\\?, updated_at=\\? WHERE ID = \\? AND status = 0"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(1, any{}, or.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	o := orderRepo.NewMysqlOrderRepository(db)

	err = o.Update(context.TODO(), or)
	assert.NoError(t, err)
	assert.NotNil(t, or)
	assert.Equal(t, or.Status, models.Taken)
	assert.NotEqual(t, now, or.UpdatedAt)
}
