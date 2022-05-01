package app

// TODO validate lat/lng
type SaveLocationRequest struct {
	VehicleId string  `json:"vehicle_id" validate:"required"`
	Lat       float64 `json:"lat" validate:"required,gte=-90,lte=90"`
	Lng       float64 `json:"lng" validate:"required,gte=-180,lte=180"`
} // @name SaveLocationRequest

type SearchLocationRequest struct {
	Lat float64 `json:"lat" validate:"required"`
	Lng float64 `json:"lng" validate:"required"`
} // @name SearchLocationRequest
