package model

import (
	"reflect"
	"testing"
)

func TestVehicle_MarshalJSON(t *testing.T) {
	type fields struct {
		Id     string
		Name   string
		Plate  string
		Type   string
		Class  string
		Seats  int
		Driver Driver
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "marshal vehicle",
			fields: fields{
				Id:    "vehicle_id",
				Name:  "name",
				Plate: "plate",
				Type:  "type",
				Class: "class",
				Seats: 1,
				Driver: Driver{
					Id:      "driver_id",
					Name:    "driver_name",
					Email:   "driver_email",
					Picture: "driver_picture",
				},
			},
			want:    []byte(`{"vehicle_id":"vehicle_id","name":"name","plate":"plate","type":"type","class":"class","seats":1,"driver":{"user_id":"driver_id","name":"driver_name","nickname":"","email":"driver_email","picture":"driver_picture"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vehicle{
				Id:     tt.fields.Id,
				Name:   tt.fields.Name,
				Plate:  tt.fields.Plate,
				Type:   tt.fields.Type,
				Class:  tt.fields.Class,
				Seats:  tt.fields.Seats,
				Driver: tt.fields.Driver,
			}
			got, err := v.MarshalJson()
			if (err != nil) != tt.wantErr {
				t.Errorf("Vehicle.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Vehicle.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicle_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Id     string
		Name   string
		Plate  string
		Type   string
		Class  string
		Seats  int
		Driver Driver
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "unmarshal vehicle",
			fields: fields{
				Id:    "vehicle_id",
				Name:  "name",
				Plate: "plate",
				Type:  "type",
				Class: "class",
				Seats: 1,
				Driver: Driver{
					Id:      "driver_id",
					Name:    "driver_name",
					Email:   "driver_email",
					Picture: "driver_picture",
				},
			},
			args: args{
				data: []byte(`{"vehicle_id":"vehicle_id","name":"name","plate":"plate","type":"type","class":"class","seats":1,"driver":{"user_id":"driver_id","name":"driver_name","nickname":"","email":"driver_email","picture":"driver_picture"}}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Vehicle{}
			if err := v.UnmarshalJson(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Vehicle.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			v2 := &Vehicle{
				Id:     tt.fields.Id,
				Name:   tt.fields.Name,
				Plate:  tt.fields.Plate,
				Type:   tt.fields.Type,
				Class:  tt.fields.Class,
				Seats:  tt.fields.Seats,
				Driver: tt.fields.Driver,
			}

			if !reflect.DeepEqual(*v, *v2) {
				t.Errorf("Vehicle.UnmarshalJSON() = %v, want %v", *v, *v2)
			}
		})
	}
}
