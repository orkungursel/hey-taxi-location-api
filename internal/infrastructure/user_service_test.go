package infrastructure

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"github.com/orkungursel/hey-taxi-location-api/internal/domain/model"
	. "github.com/orkungursel/hey-taxi-location-api/pkg/logger/mock"
	userServiceGrpc "github.com/orkungursel/hey-taxi-location-api/proto"
	userServiceGrpcMock "github.com/orkungursel/hey-taxi-location-api/proto/mock"
)

func TestUserService_GetUsersByIds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewLoggerMock()

	type args struct {
		ctx context.Context
		in  []string
	}
	tests := []struct {
		name    string
		args    args
		client  func() *userServiceGrpcMock.MockUserServiceClient
		want    map[string]model.User
		wantErr bool
	}{
		{
			name: "should return error when empty user ids",
			args: args{
				ctx: context.Background(),
				in:  []string{},
			},
			client: func() *userServiceGrpcMock.MockUserServiceClient {
				userServiceGrpcClientMock := userServiceGrpcMock.NewMockUserServiceClient(ctrl)
				return userServiceGrpcClientMock
			},
			wantErr: true,
		},
		{
			name: "should return users when success",
			args: args{
				ctx: context.Background(),
				in:  []string{"1", "2"},
			},
			client: func() *userServiceGrpcMock.MockUserServiceClient {
				userServiceGrpcClientMock := userServiceGrpcMock.NewMockUserServiceClient(ctrl)
				userServiceGrpcClientMock.EXPECT().
					GetUserInfo(gomock.Any(), gomock.Any()).
					Return(&userServiceGrpc.GetUserInfoResponse{
						Users: []*userServiceGrpc.UserInfo{
							{
								Id:     "1",
								Name:   "name1",
								Email:  "email1",
								Type:   "type1",
								Role:   "role1",
								Avatar: "avatar1",
							},
							{
								Id:     "2",
								Name:   "name2",
								Email:  "email2",
								Type:   "type2",
								Role:   "role2",
								Avatar: "avatar2",
							},
						},
					}, nil).Times(1)
				return userServiceGrpcClientMock
			},
			want: map[string]model.User{
				"1": {
					Id:       "1",
					Name:     "name1",
					Email:    "email1",
					Nickname: "name1",
					Picture:  "avatar1",
				},
				"2": {
					Id:       "2",
					Name:     "name2",
					Email:    "email2",
					Nickname: "name2",
					Picture:  "avatar2",
				},
			},
		},
		{
			name: "should return empty when not found",
			args: args{
				ctx: context.Background(),
				in:  []string{"3"},
			},
			client: func() *userServiceGrpcMock.MockUserServiceClient {
				userServiceGrpcClientMock := userServiceGrpcMock.NewMockUserServiceClient(ctrl)
				userServiceGrpcClientMock.EXPECT().GetUserInfo(gomock.Any(), gomock.Any()).
					Return(&userServiceGrpc.GetUserInfoResponse{}, nil).Times(1)
				return userServiceGrpcClientMock
			},
			want: map[string]model.User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := NewUserService(config.New(), logger, tt.client())

			got, err := us.GetUsersByIds(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetUsersByIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetUsersByIds() = %#+v, want %#+v", got, tt.want)
			}
		})
	}
}
