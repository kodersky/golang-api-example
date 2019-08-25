package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
)

type mysqlOrderRepository struct {
	Conn *sql.DB
}

// NewMysqlOrderRepository will create an object that represent the order.Repository interface
func NewMysqlOrderRepository(Conn *sql.DB) order.Repository {
	return &mysqlOrderRepository{Conn}
}

func (m *mysqlOrderRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Order, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Order, 0)
	for rows.Next() {
		t := new(models.Order)
		err = rows.Scan(
			&t.ID,
			&t.UUID,
			&t.Distance,
			&t.Status,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlOrderRepository) Fetch(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	query := `SELECT id, uuid, distance, status, updated_at, created_at
	FROM orders
	INNER JOIN (
	SELECT id
	FROM orders
	ORDER BY id
	LIMIT ? OFFSET ?)
	AS my_orders USING(id)`

	res, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return res, err
}
func (m *mysqlOrderRepository) GetByID(ctx context.Context, id string) (res *models.Order, err error) {
	query := `SELECT id, uuid, distance, status, updated_at, created_at
  						FROM orders WHERE UUID = ?`

	idB := []byte(id)
	idPB, err := uuid.ParseBytes(idB)
	if err != nil {
		return nil, err
	}
	i, err := idPB.MarshalBinary()
	if err != nil {
		return nil, err
	}

	list, err := m.fetch(ctx, query, i)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *mysqlOrderRepository) Store(ctx context.Context, o *models.Order) error {
	query := `INSERT orders SET distance=?, uuid=?, status=?, start_lat=?,
	start_long=?, end_lat=?, end_long=?, updated_at=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	uid := uuid.New()

	o.UUID = uid
	o.Status = models.Unassigned
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt

	uuidBin, err := o.UUID.MarshalBinary()

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, o.Distance, uuidBin, o.Status,
		o.StartLat, o.StartLong, o.EndLat, o.EndLong, o.UpdatedAt, o.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	o.ID = lastID
	return nil
}

func (m *mysqlOrderRepository) Update(ctx context.Context, o *models.Order) error {
	// Atomic update. No race condition.
	query := fmt.Sprintf(`UPDATE orders set status=?, updated_at=? WHERE ID = ? AND status = %d`, models.Unassigned)
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	uAt := time.Now()

	res, err := stmt.ExecContext(ctx, o.Status, uAt, o.ID)
	if err != nil {
		return err
	}

	o.UpdatedAt = uAt

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
		return err
	}

	return nil
}
