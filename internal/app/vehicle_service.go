//go:generate mockgen -source vehicle_service.go -destination mock/vehicle_service_mock.go -package mock
package app

import (
	"context"

	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
)

type VehicleService interface {
	GetVehicleById(ctx context.Context, vehicleId string) (*model.Vehicle, error)
}
