package model

import "encoding/json"

type User struct {
	Id       string   `json:"user_id"`
	Name     string   `json:"name"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Picture  string   `json:"picture"`
	Meta     metadata `json:"app_metadata"`
} // @name UserDetails

type metadata struct {
	Type string `json:"type"`
	Foo  string `json:"foo"`
}

func (ud *User) ToJson() ([]byte, error) {
	return json.Marshal(ud)
}

func (ud *User) FromJson(data []byte) error {
	return json.Unmarshal(data, &ud)
}
