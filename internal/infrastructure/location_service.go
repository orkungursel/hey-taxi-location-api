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
	ErrEmptyUserId = errors.New("user id is empty")
)

type LocationService struct {
	repo   app.LocationRepository
	logger logger.ILogger
	user   app.UserService
}

func NewLocationService(repo app.LocationRepository, logger logger.ILogger, us app.UserService) *LocationService {
	return &LocationService{
		repo:   repo,
		logger: logger,
		user:   us,
	}
}

// SaveLocation saves the location of the driver
func (s *LocationService) SaveLocation(ctx context.Context, userId string, in app.SaveLocationRequest) error {
	if userId == "" {
		return ErrEmptyUserId
	}

	l := model.Location{
		Driver: userId,
		Lat:    in.Lat,
		Lng:    in.Lng,
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

	userIds := make([]string, 0)
	for _, v := range res {
		userIds = append(userIds, v.Driver)
	}

	if len(userIds) == 0 {
		return data, nil
	}

	users, err := s.user.GetUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		u, ok := users[v.Driver]
		if !ok {
			s.logger.Warnf(`User "%s" not found`, v.Driver)
			continue
		}

		data = append(data, app.LocationResponse{
			Driver: *MapUserToDriver(u),
			Lat:    v.Lat,
			Lng:    v.Lng,
			Dist:   v.Dist,
		})
	}

	return data, err
}
