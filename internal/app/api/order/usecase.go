package order

import (
	"context"

	"github.com/kodersky/golang-api-example/internal/app/api/models"
)

// Usecase represent the order's usecases
type Usecase interface {
	Fetch(ctx context.Context, limit int, offset int) ([]*models.Order, error)
	GetByID(ctx context.Context, id string) (*models.Order, error)
	Update(ctx context.Context, or *models.Order) error
	Store(context.Context, *models.Order) error
}
