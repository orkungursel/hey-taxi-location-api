package infrastructure

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger"
)

const (
	dbKey        = "drivers" // dbKey is the key to store drivers in redis
	maxLimit     = 100       // maxLimit is the maximum limit for the search
	defaultLimit = 20        // defaultLimit is the default limit for the search
)

type LocationRepository struct {
	db     *redis.Client
	logger logger.ILogger
	dbKey  string
}

func NewLocationRepository(db *redis.Client, logger logger.ILogger) *LocationRepository {
	return &LocationRepository{
		db:     db,
		logger: logger,
		dbKey:  dbKey,
	}
}

// Save saves the location of the driver to redis database
func (r *LocationRepository) Save(ctx context.Context, in model.Location) error {
	d := MapLocationToRedisGeoLocation(in)
	return r.db.GeoAdd(ctx, r.dbKey, d).Err()
}

// Search searches for drivers in redis database
func (r *LocationRepository) Search(ctx context.Context, lat, lng, radius float64,
	unit string, limit int) ([]model.Location, error) {

	if limit == 0 || limit > maxLimit {
		limit = defaultLimit
	}

	q := &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      unit,
		WithCoord: true,
		WithDist:  true,
		Count:     limit,
	}

	d, err := r.db.GeoRadius(ctx, r.dbKey, lat, lng, q).Result()

	if err != nil {
		return nil, err
	}

	var res []model.Location = make([]model.Location, len(d))
	for i, v := range d {
		res[i] = *MapRedisGeoLocationToDomain(v)
	}

	return res, nil
}
