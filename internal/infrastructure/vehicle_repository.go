package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

const (
	vehicleDbKey        = "vehicle"                     // dbKey is the key to store drivers in redis
	vehicleDbExpiration = time.Duration(24) * time.Hour // expire is the expiration time for the vehicle in redis
)

type VehicleRepository struct {
	db     *redis.Client
	logger logger.ILogger
	dbKey  string
	expire time.Duration
}

func NewVehicleRepository(db *redis.Client, logger logger.ILogger) *VehicleRepository {
	return &VehicleRepository{
		db:     db,
		logger: logger,
		dbKey:  vehicleDbKey,
		expire: vehicleDbExpiration,
	}
}

// generateDbKey generates the key to store the vehicle in redis
func (r *VehicleRepository) generateDbKey(vehicleId string) (string, error) {
	if vehicleId == "" {
		return "", errors.New("vehicleId is empty")
	}

	return r.dbKey + ":" + vehicleId, nil
}

// Get returns the vehicle from redis database
func (r *VehicleRepository) Get(ctx context.Context, vehicleId string) (*model.Vehicle, error) {
	key, err := r.generateDbKey(vehicleId)
	if err != nil {
		return nil, err
	}

	s, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	vehicle := &model.Vehicle{}
	if err := vehicle.UnmarshalJson([]byte(s)); err != nil {
		return nil, err
	}

	return vehicle, nil
}

// Save saves the vehicle to redis database
func (r *VehicleRepository) Save(ctx context.Context, vehicle *model.Vehicle) error {
	key, err := r.generateDbKey(vehicle.Id)
	if err != nil {
		return err
	}

	s, err := vehicle.MarshalJson()
	if err != nil {
		return err
	}

	return r.db.SetEX(ctx, key, s, r.expire).Err()
}

// Delete deletes the vehicle from redis database
func (r *VehicleRepository) Delete(ctx context.Context, vehicleId string) error {
	key, err := r.generateDbKey(vehicleId)
	if err != nil {
		return err
	}

	return r.db.Del(ctx, key).Err()
}
