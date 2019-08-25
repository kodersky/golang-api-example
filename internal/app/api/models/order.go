package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus uint8

func (os OrderStatus) String() string {
	return [...]string{"UNASSIGNED", "TAKEN"}[os]
}

func (os OrderStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + os.String() + `"`), nil
}

// Available order statuses
const (
	Unassigned OrderStatus = iota
	Taken
)

// Order represent the order model
type Order struct {
	ID        int64       `json:"-"`
	UUID      uuid.UUID   `json:"id"`
	Status    OrderStatus `json:"status"`
	Distance  int         `json:"distance"`
	StartLat  float64     `json:"-"`
	StartLong float64     `json:"-"`
	EndLat    float64     `json:"-"`
	EndLong   float64     `json:"-"`
	UpdatedAt time.Time   `json:"-"`
	CreatedAt time.Time   `json:"-"`
}
