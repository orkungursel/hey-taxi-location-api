//go:generate mockgen -source location_repository.go -destination mock/location_repository_mock.go -package mock
package app

import (
	"context"

	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
)

type LocationRepository interface {
	Save(ctx context.Context, location model.Location) error
	Search(
		ctx context.Context,
		lat, lng, radius float64,
		unit string,
		limit int,
	) ([]model.Location, error)
}
