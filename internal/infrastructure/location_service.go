package infrastructure

import (
	"context"
	"errors"

	"github.com/orkungursel/hey-taxi-location-api/internal/app"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

const (
	defaultLocationSearchRadius = 200
)

var (
	ErrEmptyUserId          = errors.New("user id is empty")
	ErrVehicleService       = errors.New("vehicle service error")
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrVehicleOwnerNotMatch = errors.New("vehicle owner not match")
)

type LocationService struct {
	repo           app.LocationRepository
	vehicleService app.VehicleService
	logger         logger.ILogger
}

func NewLocationService(repo app.LocationRepository, logger logger.ILogger,
	vehicleService app.VehicleService) *LocationService {
	return &LocationService{
		repo:           repo,
		logger:         logger,
		vehicleService: vehicleService,
	}
}

// SaveLocation saves the location of the driver
func (s *LocationService) SaveLocation(ctx context.Context, userId string, in app.SaveLocationRequest) error {
	if userId == "" {
		return ErrEmptyUserId
	}

	if err := app.Validate(in); err != nil {
		return err
	}

	vehicle, err := s.vehicleService.GetVehicleById(ctx, in.VehicleId)
	if err != nil {
		s.logger.Error(ctx, "vehicle service error", err)
		return ErrVehicleService
	}

	if vehicle == nil {
		return ErrVehicleNotFound
	}

	if vehicle.Driver.Id != userId {
		s.logger.Info(ctx, "vehicle owner not match", vehicle, userId)
		return ErrVehicleOwnerNotMatch
	}

	l := model.Location{
		VehicleId: in.VehicleId,
		Lat:       in.Lat,
		Lng:       in.Lng,
	}

	if err := app.Validate(l); err != nil {
		return err
	}

	return s.repo.Save(ctx, l)
}

// Search searches for drivers
func (s *LocationService) SearchLocations(ctx context.Context, q app.SearchLocationRequest) ([]app.LocationResponse, error) {
	res, err := s.repo.Search(ctx, q.Lat, q.Lng, defaultLocationSearchRadius, "km", 0)
	if err != nil {
		return nil, err
	}

	data := make([]app.LocationResponse, 0)

	for _, v := range res {
		vehicle, err := s.vehicleService.GetVehicleById(ctx, v.VehicleId)
		if err != nil {
			s.logger.Error(ctx, "vehicle service error", err)
			return nil, ErrVehicleService
		}

		if vehicle == nil {
			continue
		}

		data = append(data, app.LocationResponse{
			Vehicle: *vehicle,
			Lat:     v.Lat,
			Lng:     v.Lng,
			Dist:    v.Dist,
		})
	}

	return data, err
}
