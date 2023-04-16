package models

type User struct {
	Email       string `json:"email"`
	Fullname    string `json:"fullName"`
	Iat         int    `json:"iat"`
	Id          int    `json:"id"`
	IsVerified  bool   `json:"isVerified"`
	PhoneNumber string `json:"phoneNumber"`
	Username    string `json:"username"`
}