package app

import "github.com/orkungursel/hey-taxi-location-api/internal/domain/model"

type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
} // @name HTTPError

type LocationResponse struct {
	Vehicle model.Vehicle `json:"vehicle"`
	Lat     float64       `json:"lat"`
	Lng     float64       `json:"lng"`
	Dist    float64       `json:"dist"`
} // @name LocationResponse
