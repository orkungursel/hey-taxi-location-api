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
	logger "github.com/orkungursel/hey-taxi-location-api/pkg/logger/mock"
)

var (
	u1 = model.User{Id: "driver1", Name: "Driver Name", Nickname: "Driver Nickname", Email: "driver@example.com", Picture: "driver.png"}
	u2 = model.User{Id: "driver2", Name: "Driver Name 2", Nickname: "Driver Nickname 2", Email: "driver2@example.com", Picture: "driver2.png"}
	d1 = MapUserToDriver(u1)
	d2 = MapUserToDriver(u2)
	v1 = model.Vehicle{Id: "vehicle1", Name: "Vehicle Name", Plate: "plate", Type: "type", Class: "class", Seats: 4, Driver: *d1}
	v2 = model.Vehicle{Id: "vehicle2", Name: "Vehicle Name 2", Plate: "plate 2", Type: "type 2", Class: "class 2", Seats: 5, Driver: *d2}
	l1 = model.Location{VehicleId: v1.Id, Lat: 1.0, Lng: 1.0}
	l2 = model.Location{VehicleId: v2.Id, Lat: 20.0, Lng: 20.0}
)

func TestLocationService_SaveLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewLoggerMock()

	type args struct {
		userId string
		in     app.SaveLocationRequest
	}
	tests := []struct {
		name           string
		args           args
		repository     func() app.LocationRepository
		vehicleService func() app.VehicleService
		wantErr        bool
	}{
		{
			name: "should success when data is valid",
			repository: func() app.LocationRepository {
				r := mock.NewMockLocationRepository(ctrl)
				r.EXPECT().Save(gomock.Any(), model.Location{VehicleId: v1.Id, Lat: 1.0, Lng: 1.0}).Return(nil).Times(1)
				return r
			},
			vehicleService: func() app.VehicleService {
				s := mock.NewMockVehicleService(ctrl)
				s.EXPECT().GetVehicleById(gomock.Any(), v1.Id).Return(&v1, nil).Times(1)
				return s
			},
			args: args{
				userId: d1.Id,
				in: app.SaveLocationRequest{
					VehicleId: v1.Id,
					Lat:       1.0,
					Lng:       1.0,
				},
			},
		},
		{
			name: "should fail when no data provided",
			repository: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
				return userService
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "should fail when user id is empty",
			repository: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
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
			repository: func() app.LocationRepository {
				repo := mock.NewMockLocationRepository(ctrl)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
				return userService
			},
			args: args{
				userId: v1.Id,
				in: app.SaveLocationRequest{
					Lat: 300.0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locationService := NewLocationService(tt.repository(), loggerMock, tt.vehicleService())

			if err := locationService.SaveLocation(context.Background(), tt.args.userId, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("LocationService.SaveLocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocationService_SearchLocations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := logger.NewLoggerMock()

	type args struct {
		q app.SearchLocationRequest
	}
	tests := []struct {
		name           string
		args           args
		repository     func() app.LocationRepository
		vehicleService func() app.VehicleService
		want           []app.LocationResponse
		wantErr        bool
	}{
		{
			name: "should return locations when data is valid",
			repository: func() app.LocationRepository {
				repo, _ := SetupLocationRepositoryMocks()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			vehicleService: func() app.VehicleService {
				vs := mock.NewMockVehicleService(ctrl)
				vs.EXPECT().GetVehicleById(gomock.Any(), v1.Id).
					Return(&v1, nil).Times(1)
				return vs
			},
			args: args{
				q: app.SearchLocationRequest{
					Lat: 1.0,
					Lng: 1.0,
				},
			},
			want: []app.LocationResponse{
				{Vehicle: v1, Lat: 1.0, Lng: 1.0, Dist: 0.0001},
			},
		},
		{
			name: "should skip location when vehicle is not found",
			repository: func() app.LocationRepository {
				repo, _ := SetupLocationRepositoryMocks()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
				userService.EXPECT().GetVehicleById(gomock.Any(), gomock.Any()).
					Return(nil, nil).Times(1)
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
			repository: func() app.LocationRepository {
				repo, _ := SetupLocationRepositoryMocks()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
				userService.EXPECT().GetVehicleById(gomock.Any(), gomock.Any()).Times(0)
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
			repository: func() app.LocationRepository {
				repo, redis := SetupLocationRepositoryMocks()
				redis.Close()
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
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
			name: "should return error when vehicle service fails",
			repository: func() app.LocationRepository {
				repo, _ := SetupLocationRepositoryMocks()
				_ = repo.Save(context.Background(), l1)
				_ = repo.Save(context.Background(), l2)
				return repo
			},
			vehicleService: func() app.VehicleService {
				userService := mock.NewMockVehicleService(ctrl)
				userService.EXPECT().GetVehicleById(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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
			locationService := NewLocationService(tt.repository(), loggerMock, tt.vehicleService())
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
