package infrastructure

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/orkungursel/hey-taxi-location-api/internal/app"
	"github.com/orkungursel/hey-taxi-location-api/internal/app/mock"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	. "github.com/orkungursel/hey-taxi-location-api/pkg/logger/mock"
)

var (
	u1 = model.User{Id: "driver1", Name: "Driver Name", Nickname: "Driver Nickname", Email: "driver@example.com", Picture: "driver.png"}
	u2 = model.User{Id: "driver2", Name: "Driver Name 2", Nickname: "Driver Nickname 2", Email: "driver2@example.com", Picture: "driver2.png"}
	d1 = MapUserToDriver(u1)
	d2 = MapUserToDriver(u2)
	l1 = model.Location{Driver: u1.Id, Lat: 1.0, Lng: 1.0}
	l2 = model.Location{Driver: u2.Id, Lat: 20.0, Lng: 20.0}
)

func TestLocationService_SaveLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewLoggerMock()

	type args struct {
		userId string
		in     app.SaveLocationRequest
	}
	tests := []struct {
		name        string
		args        args
		repo        func() app.LocationRepository
		userService func() app.UserService
		wantErr     bool
	}{
		{
			name: "should success when data is valid",
			repo: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				repo.EXPECT().Save(gomock.Any(), model.Location{Driver: u1.Id, Lat: 1.0, Lng: 1.0}).Return(nil).Times(1)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				return userService
			},
			args: args{
				userId: u1.Id,
				in: app.SaveLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
		},
		{
			name: "should fail when no data provided",
			repo: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				return userService
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "should fail when user id is empty",
			repo: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				return userService
			},
			args: args{
				in: app.SaveLocationRequest{
					Lat: 100.0,
					Lng: 100.0,
				},
			},
			wantErr: true,
		},
		{
			name: "should fail when add location data is invalid",
			repo: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				return userService
			},
			args: args{
				userId: u1.Id,
				in: app.SaveLocationRequest{
					Lat: 300.0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locationService := NewLocationService(tt.repo(), loggerMock, tt.userService())

			if err := locationService.SaveLocation(context.Background(), tt.args.userId, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("LocationService.SaveLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocationService_SearchLocations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := NewLoggerMock()

	type args struct {
		q app.SearchLocationRequest
	}
	tests := []struct {
		name        string
		args        args
		repo        func() app.LocationRepository
		userService func() app.UserService
		want        []app.LocationResponse
		wantErr     bool
	}{
		{
			name: "should return locations when data is valid",
			repo: func() app.LocationRepository {
				repo, _ := SetupRepositoryMock()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				userService.EXPECT().GetUsersByIds(gomock.Any(), gomock.Any()).
					Return(map[string]model.User{u1.Id: u1, u2.Id: u2}, nil).Times(1)
				return userService
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
			want: []app.LocationResponse{
				{Driver: *d1, Lat: 1.0, Lng: 1.0, Dist: 0.0001},
			},
		},
		{
			name: "should skip location when user is not found",
			repo: func() app.LocationRepository {
				repo, _ := SetupRepositoryMock()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				userService.EXPECT().GetUsersByIds(gomock.Any(), gomock.Any()).
					Return(map[string]model.User{}, nil).Times(1)
				return userService
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
			want: []app.LocationResponse{},
		},
		{
			name: "should return empty list when no data found",
			repo: func() app.LocationRepository {
				repo, _ := SetupRepositoryMock()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				userService.EXPECT().GetUsersByIds(gomock.Any(), gomock.Any()).AnyTimes()
				return userService
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 30.0,
					Lng: 30.0,
				},
			},
			want: []app.LocationResponse{},
		},
		{
			name: "should return error when user repository fails",
			repo: func() app.LocationRepository {
				repo, redis := SetupRepositoryMock()
				redis.Close()
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				return userService
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when user service fails",
			repo: func() app.LocationRepository {
				repo, _ := SetupRepositoryMock()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			userService: func() app.UserService {
				userService := mock.NewMockUserService(ctrl)
				userService.EXPECT().GetUsersByIds(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				return userService
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locationService := NewLocationService(tt.repo(), loggerMock, tt.userService())
			got, err := locationService.SearchLocations(context.Background(), tt.args.q)

			if (err != nil) != tt.wantErr {
				t.Errorf("LocationService.SearchLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocationService.SearchLocations() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
