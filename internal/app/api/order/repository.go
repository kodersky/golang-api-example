package order

import (
	"context"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
)

// Repository represent the order's repository contract
type Repository interface {
	Fetch(ctx context.Context, limit int, offset int) (res []*models.Order, err error)
	Update(ctx context.Context, o *models.Order) error
	Store(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
}
