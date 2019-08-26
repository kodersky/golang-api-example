package tests

import (
	"context"
	"errors"
	"github.com/kodersky/golang-api-example/internal/app/api/order"
	"googlemaps.github.io/maps"
	"time"
)

type FakeClient struct {
	Gr maps.DistanceMatrixResponse
	T  time.Duration
	E  error
}

var GoogleResponse = maps.DistanceMatrixResponse{
	OriginAddresses:      []string{"1119 Rama IX Rd, Khwaeng Suan Luang, Khet Suan Luang, Krung Thep Maha Nakhon 10250, Thailand"},
	DestinationAddresses: []string{"2098/257 Soi Ramkhamhaeng 24 Yeak 30, Khwaeng Hua Mak, Khet Bang Kapi, Krung Thep Maha Nakhon 10240, Thailand"},
	Rows: []maps.DistanceMatrixElementsRow{
		{
			Elements: []*maps.DistanceMatrixElement{
				{
					Status:            "OK",
					Duration:          212000000000,
					DurationInTraffic: 0,
					Distance: maps.Distance{
						HumanReadable: "1.4 km",
						Meters:        1426,
					},
				},
			},
		},
	},
}

func (fk *FakeClient) DistanceMatrix(ctx context.Context, r *maps.DistanceMatrixRequest) (*maps.DistanceMatrixResponse, error) {
	if fk.T > 0 {
		time.Sleep(4 * time.Second)
	}
	if fk.E != nil {
		return &maps.DistanceMatrixResponse{}, errors.New("google api problem")
	}
	return &GoogleResponse, nil
}

func NewWithClient(gr maps.DistanceMatrixResponse, t time.Duration, e error) *order.Client {
	var c FakeClient
	c.Gr = gr
	c.T = t
	c.E = e
	return &order.Client{
		Client: &c,
	}
}
