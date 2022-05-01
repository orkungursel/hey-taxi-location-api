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
	"github.com/orkungursel/hey-taxi-location-api/proto"
	protoMock "github.com/orkungursel/hey-taxi-location-api/proto/mock"
)

var vehicle1 = &model.Vehicle{
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

func TestVehicleService_GetVehicleById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		repo   func() app.VehicleRepository
		client func() proto.VehicleServiceClient
	}
	type args struct {
		ctx       context.Context
		vehicleId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Vehicle
		wantErr bool
	}{
		{
			name: "should return vehicle from repository",
			fields: fields{
				repo: func() app.VehicleRepository {
					repo := mock.NewMockVehicleRepository(ctrl)
					repo.EXPECT().Get(gomock.Any(), vehicle1.Id).DoAndReturn(
						func(_ context.Context, vehicleId string) (*model.Vehicle, error) {
							if vehicleId != vehicle1.Id {
								return nil, nil
							}
							return vehicle1, nil
						},
					)
					return repo
				},
				client: func() proto.VehicleServiceClient {
					client := protoMock.NewMockVehicleServiceClient(ctrl)
					client.EXPECT().GetVehicle(gomock.Any(), gomock.Any()).Times(0)
					return client
				},
			},
			args: args{
				ctx:       context.Background(),
				vehicleId: vehicle1.Id,
			},
			want: vehicle1,
		},
		{
			name: "should return vehicle from grpc service when repository returns nil",
			fields: fields{
				repo: func() app.VehicleRepository {
					repo := mock.NewMockVehicleRepository(ctrl)
					repo.EXPECT().Get(gomock.Any(), vehicle1.Id).Return(nil, nil).Times(1)
					repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(1)
					return repo
				},
				client: func() proto.VehicleServiceClient {
					client := protoMock.NewMockVehicleServiceClient(ctrl)
					client.EXPECT().GetVehicle(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
						func(_ context.Context, req *proto.GetVehicleRequest, _ ...interface{}) (*proto.GetVehicleResponse, error) {
							if req.Id != vehicle1.Id {
								return nil, nil
							}
							return &proto.GetVehicleResponse{
								Id:    vehicle1.Id,
								Name:  vehicle1.Name,
								Plate: vehicle1.Plate,
								Type:  vehicle1.Type,
								Class: vehicle1.Class,
								Seats: int32(vehicle1.Seats),
								Driver: &proto.DriverDetailsResponse{
									Id:     vehicle1.Driver.Id,
									Name:   vehicle1.Driver.Name,
									Email:  vehicle1.Driver.Email,
									Avatar: vehicle1.Driver.Picture,
								},
							}, nil
						},
					)
					return client
				},
			},
			args: args{
				ctx:       context.Background(),
				vehicleId: vehicle1.Id,
			},
			want: vehicle1,
		},
		{
			name: "should return error when grpc service returns error",
			fields: fields{
				repo: func() app.VehicleRepository {
					repo := mock.NewMockVehicleRepository(ctrl)
					repo.EXPECT().Get(gomock.Any(), vehicle1.Id).Return(nil, nil).Times(1)
					repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).Times(0)
					return repo
				},
				client: func() proto.VehicleServiceClient {
					client := protoMock.NewMockVehicleServiceClient(ctrl)
					client.EXPECT().GetVehicle(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
					return client
				},
			},
			args: args{
				ctx:       context.Background(),
				vehicleId: vehicle1.Id,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := &VehicleService{
				repo:   tt.fields.repo(),
				client: tt.fields.client(),
				logger: logger.NewLoggerMock(),
			}
			got, err := vs.GetVehicleById(tt.args.ctx, tt.args.vehicleId)
			if (err != nil) != tt.wantErr {
				t.Errorf("VehicleService.GetVehicleById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VehicleService.GetVehicleById() = %v, want %v", got, tt.want)
			}
		})
	}
}
