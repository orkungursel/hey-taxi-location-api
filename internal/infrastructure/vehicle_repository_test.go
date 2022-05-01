package infrastructure

import (
	"context"
	"reflect"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	"github.com/orkungursel/hey-taxi-location-api/pkg/logger/mock"
)

func SetupVehicleRepositoryMocks() (*VehicleRepository, *redis.Client) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	r := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	repo := NewVehicleRepository(r, mock.NewLoggerMock())
	return repo, r
}

func TestVehicleRepository_generateDbKey(t *testing.T) {
	repo, _ := SetupVehicleRepositoryMocks()

	type args struct {
		vehicleId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		dbKey   string
	}{
		{
			name: "generate key",
			args: args{
				vehicleId: "vehicle_id",
			},
			want: vehicleDbKey + ":vehicle_id",
		},
		{
			name: "generate key with custom db key",
			args: args{
				vehicleId: "vehicle_id",
			},
			want:  "custom_key:vehicle_id",
			dbKey: "custom_key",
		},
		{
			name: "should return error when vehicle id is empty",
			args: args{
				vehicleId: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbKey != "" {
				repo.dbKey = tt.dbKey
			}

			if got, err := repo.generateDbKey(tt.args.vehicleId); got != tt.want {
				if (err != nil) != tt.wantErr {
					t.Errorf("VehicleRepository.generateDbKey() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				t.Errorf("VehicleRepository.generateDbKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicleRepository_Get(t *testing.T) {
	repo, _ := SetupVehicleRepositoryMocks()
	ctx := context.Background()

	v1 := &model.Vehicle{
		Id:    "vehicle_id",
		Name:  "name",
		Plate: "plate",
		Type:  "type",
		Class: "class",
		Seats: 1,
		Driver: model.Driver{
			Id:      "driver_id",
			Name:    "driver_name",
			Email:   "driver_email",
			Picture: "driver_picture",
		},
	}

	v2 := &model.Vehicle{
		Id:    "vehicle_id_2",
		Name:  "name",
		Plate: "plate",
		Type:  "type",
		Class: "class",
		Seats: 1,
		Driver: model.Driver{
			Id:      "driver_id",
			Name:    "driver_name",
			Email:   "driver_email",
			Picture: "driver_picture",
		},
	}

	repo.Save(ctx, v1)
	repo.Save(ctx, v2)

	type args struct {
		ctx       context.Context
		vehicleId string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Vehicle
		wantErr bool
	}{
		{
			name: "get vehicle",
			args: args{
				ctx:       ctx,
				vehicleId: v1.Id,
			},
			want:    v1,
			wantErr: false,
		},
		{
			name: "get vehicle with invalid id",
			args: args{
				ctx:       ctx,
				vehicleId: "invalid_id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(tt.args.ctx, tt.args.vehicleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("VehicleRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VehicleRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicleRepository_Save(t *testing.T) {
	repo, _ := SetupVehicleRepositoryMocks()
	ctx := context.Background()

	vehicle := &model.Vehicle{
		Id:    "vehicle_id",
		Name:  "name",
		Plate: "plate",
		Type:  "type",
		Class: "class",
		Seats: 1,
		Driver: model.Driver{
			Id:      "driver_id",
			Name:    "driver_name",
			Email:   "driver_email",
			Picture: "driver_picture",
		},
	}

	vehicle_invalid := &model.Vehicle{
		Id: "",
	}

	type args struct {
		ctx     context.Context
		vehicle *model.Vehicle
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save vehicle",
			args: args{
				ctx:     ctx,
				vehicle: vehicle,
			},
		},
		{
			name: "save vehicle with invalid id",
			args: args{
				ctx:     ctx,
				vehicle: vehicle_invalid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Save(tt.args.ctx, tt.args.vehicle); (err != nil) != tt.wantErr {
				t.Errorf("VehicleRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVehicleRepository_Delete(t *testing.T) {
	repo, _ := SetupVehicleRepositoryMocks()
	ctx := context.Background()

	v1 := &model.Vehicle{
		Id:    "vehicle_id",
		Name:  "name",
		Plate: "plate",
		Type:  "type",
		Class: "class",
		Seats: 1,
		Driver: model.Driver{
			Id:      "driver_id",
			Name:    "driver_name",
			Email:   "driver_email",
			Picture: "driver_picture",
		},
	}

	v2 := &model.Vehicle{
		Id:    "vehicle_id_2",
		Name:  "name",
		Plate: "plate",
		Type:  "type",
		Class: "class",
		Seats: 1,
		Driver: model.Driver{
			Id:      "driver_id",
			Name:    "driver_name",
			Email:   "driver_email",
			Picture: "driver_picture",
		},
	}

	repo.Save(ctx, v1)
	repo.Save(ctx, v2)

	type args struct {
		ctx       context.Context
		vehicleId string
	}
	tests := []struct {
		name       string
		args       args
		checkAfter func(t *testing.T, repo *VehicleRepository)
		wantErr    bool
	}{
		{
			name: "delete vehicle",
			args: args{
				ctx:       ctx,
				vehicleId: v1.Id,
			},
			checkAfter: func(t *testing.T, repo *VehicleRepository) {
				if v, err := repo.Get(ctx, v1.Id); err != nil {
					if v != nil {
						t.Errorf("VehicleRepository.Delete() item not deleted %v", v)
					}
				}

				if v, err := repo.Get(ctx, v2.Id); err != nil {
					if v == nil {
						t.Errorf("VehicleRepository.Delete() item not exists: %v", v2)
					}
				}
			},
		},
		{
			name: "delete vehicle with invalid id",
			args: args{
				ctx:       ctx,
				vehicleId: "invalid_id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Delete(tt.args.ctx, tt.args.vehicleId); (err != nil) != tt.wantErr {
				t.Errorf("VehicleRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkAfter != nil {
				tt.checkAfter(t, repo)
			}
		})
	}
}
