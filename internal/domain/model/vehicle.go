package model

import "encoding/json"

type Vehicle struct {
	Id     string `json:"vehicle_id"`
	Name   string `json:"name"`
	Plate  string `json:"plate"`
	Type   string `json:"type"`
	Class  string `json:"class"`
	Seats  int    `json:"seats"`
	Driver Driver `json:"driver"`
} // @name Vehicle

func (v *Vehicle) MarshalJson() ([]byte, error) {
	return json.Marshal(*v)
}

func (v *Vehicle) UnmarshalJson(data []byte) error {
	return json.Unmarshal(data, &v)
}
