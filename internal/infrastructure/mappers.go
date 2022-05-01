package infrastructure

import (
	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
)

func MapRedisGeoLocationToDomain(g redis.GeoLocation) *model.Location {
	return &model.Location{
		VehicleId: g.Name,
		Lat:       g.Latitude,
		Lng:       g.Longitude,
		Dist:      g.Dist,
	}
}

func MapLocationToRedisGeoLocation(l model.Location) *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      l.VehicleId,
		Longitude: l.Lng,
		Latitude:  l.Lat,
		Dist:      l.Dist,
	}
}

func MapUserToDriver(u model.User) *model.Driver {
	return &model.Driver{
		Id:       u.Id,
		Name:     u.Name,
		Nickname: u.Nickname,
		Email:    u.Email,
		Picture:  u.Picture,
	}
}
