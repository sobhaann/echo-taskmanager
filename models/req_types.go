package models

type LoginReq struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type SignupReq struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}
