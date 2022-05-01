package infrastructure

import (
	"context"
	"errors"

	"github.com/orkungursel/hey-taxi-location-api/internal/app"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
	"github.com/orkungursel/hey-taxi-location-api/proto"
)

type VehicleService struct {
	repo   app.VehicleRepository
	client proto.VehicleServiceClient
	logger logger.ILogger
}

func NewVehicleService(logger logger.ILogger,
	client proto.VehicleServiceClient, repo app.VehicleRepository) *VehicleService {
	return &VehicleService{
		repo:   repo,
		client: client,
		logger: logger,
	}
}

func (vs *VehicleService) GetVehicleById(ctx context.Context,
	vehicleId string) (*model.Vehicle, error) {

	if vehicleId == "" {
		return nil, errors.New("empty vehicle id")
	}

	if vehicle, err := vs.repo.Get(ctx, vehicleId); err == nil && vehicle != nil {
		return vehicle, nil
	}

	vehicle, err := vs.getVehicleByIdFromGrpcService(ctx, vehicleId)
	if err != nil {
		return nil, err
	}

	if err := vs.repo.Save(ctx, vehicle); err != nil {
		vs.logger.Infof("failed to save vehicle to repository: %v", err)
	}

	return vehicle, nil
}

func (vs *VehicleService) getVehicleByIdFromGrpcService(ctx context.Context,
	vehicleId string) (*model.Vehicle, error) {

	vehicle, err := vs.client.GetVehicle(ctx, &proto.GetVehicleRequest{Id: vehicleId})
	if err != nil {
		return nil, err
	}

	return &model.Vehicle{
		Id:    vehicle.Id,
		Name:  vehicle.Name,
		Plate: vehicle.Plate,
		Type:  vehicle.Type,
		Class: vehicle.Class,
		Seats: int(vehicle.Seats),
		Driver: model.Driver{
			Id:      vehicle.Driver.Id,
			Name:    vehicle.Driver.Name,
			Email:   vehicle.Driver.Email,
			Picture: vehicle.Driver.Avatar,
		},
	}, nil
}
