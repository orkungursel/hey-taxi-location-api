//go:generate mockgen -source location_service.go -destination mock/location_service_mock.go -package mock
package app

import "context"

type LocationService interface {
	SaveLocation(ctx context.Context, userId string, in SaveLocationRequest) error
	SearchLocations(ctx context.Context, req SearchLocationRequest) ([]LocationResponse, error)
}
