package model

type Driver struct {
	Id       string `json:"user_id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
} // @name Driver
