package order

import (
	"context"
	"googlemaps.github.io/maps"
)

type Client struct {
	Client GoogleMapClient
}

// GoogleMapClient represents maps.Client
type GoogleMapClient interface {
	DistanceMatrix(ctx context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error)
}
