package infrastructure

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	. "github.com/orkungursel/hey-taxi-location-api/pkg/logger/mock"
)

func newRedisClientMock() (r *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	r = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return
}

func SetupRepositoryMock() (repo *LocationRepository, redisClient *redis.Client) {
	logger := NewLoggerMock()
	redisClient = newRedisClientMock()
	repo = NewLocationRepository(redisClient, logger)
	return
}

func TestLocationRepository_Save(t *testing.T) {
	t.Parallel()

	repo, _ := SetupRepositoryMock()

	type args struct {
		in model.Location
	}

	tests := []struct {
		name    string
		r       *LocationRepository
		args    args
		wantErr bool
	}{
		{
			name: "should success",
			r:    repo,
			args: args{
				in: model.Location{
					Driver: "driver",
					Lat:    1.0,
					Lng:    1.0,
				},
			},
			wantErr: false,
		},
		{
			name: "should success with lat 0 and lng 0",
			r:    repo,
			args: args{
				in: model.Location{
					Driver: "driver-2",
					Lat:    0,
					Lng:    0,
				},
			},
			wantErr: false,
		},
		{
			name: "fail if driver is empty",
			r:    repo,
			args: args{
				in: model.Location{
					Lat: 100.0,
					Lng: 100.0,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Save(context.Background(), tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("LocationRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLocationRepository_Search(t *testing.T) {
	t.Parallel()

	repo, redis := SetupRepositoryMock()

	d1 := model.Location{Driver: "driver1", Lat: 1.0, Lng: 1.0}
	d2 := model.Location{Driver: "driver2", Lat: 20.0, Lng: 20.0}

	_ = repo.Save(context.Background(), d1)
	_ = repo.Save(context.Background(), d2)

	type args struct {
		lat    float64
		lng    float64
		radius float64
		unit   string
		limit  int
	}

	tests := []struct {
		name            string
		r               *LocationRepository
		args            args
		want            []model.Location
		wantErr         bool
		closeConnection bool
	}{
		{
			name: "should return 1 result with radius 10 km and lat 1 and lng 1",
			r:    repo,
			args: args{
				lat:    1.0,
				lng:    1.0,
				radius: 10.0,
				unit:   "km",
			},
			want: []model.Location{d1},
		},
		{
			name: "should return 2 results with radius 3000 km and lat 1 and lng 1",
			r:    repo,
			args: args{
				lat:    1.0,
				lng:    1.0,
				radius: 3000.0,
				unit:   "km",
			},
			want: []model.Location{d1, d2},
		},
		{
			name: "should return error if redis connection is closed",
			r:    repo,
			args: args{
				lat:    1.0,
				lng:    1.0,
				radius: 10.0,
				unit:   "km",
			},
			wantErr:         true,
			closeConnection: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closeConnection {
				redis.Close()
			}

			got, err := tt.r.Search(context.Background(), tt.args.lat, tt.args.lng, tt.args.radius, tt.args.unit, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocationRepository.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("LocationRepository.Search() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
