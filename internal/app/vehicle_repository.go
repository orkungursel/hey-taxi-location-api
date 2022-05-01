//go:generate mockgen -source vehicle_repository.go -destination mock/vehicle_repository_mock.go -package mock
package app

import (
	"context"

	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
)

type VehicleRepository interface {
	Get(ctx context.Context, vehicleId string) (*model.Vehicle, error)
	Save(ctx context.Context, vehicle *model.Vehicle) error
	Delete(ctx context.Context, vehicleId string) error
}
