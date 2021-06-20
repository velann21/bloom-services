package requests

import (
	"encoding/json"
	"github.com/velann21/bloom-services/common-lib/helpers"
	"io"
	"time"
)

type User struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Address      Address `json:"address"`
	DOB          DOB     `json:"dob"`
}

type DOB struct {
	Year int `json:"year"`
	Month int `json:"month"`
	Day int `json:"day"`
}

type Address struct {
	ZipCode string `json:"zip_code"`
	StreetName string `json:"street_name"`
	HouseNumber int `json:"house_number"`
}

func NewUserEntity()*User {
	return &User{}
}

func (user *User) ValidateUser() error {
	if user.Name == ""{
		return helpers.InvalidRequest
	}
	if user.Address.HouseNumber <= 0{
		return helpers.InvalidRequest
	}
	if user.Address.StreetName == ""{
		return helpers.InvalidRequest
	}
	if user.Address.ZipCode == ""{
		return helpers.InvalidRequest
	}
	if user.DOB.Day <= 0 || user.DOB.Day >= 31{
		return helpers.InvalidRequest
	}
	if user.DOB.Year >= time.Now().Year() || user.DOB.Year <= 0{
		return helpers.InvalidRequest
	}
	if user.DOB.Month <= 0 || user.DOB.Month > 12{
		return helpers.InvalidRequest
	}
	return nil
}

func (user *User) PopulateUser(body io.Reader) error {
	decode := json.NewDecoder(body)
	err := decode.Decode(user)
	if err != nil {
		return helpers.InvalidRequest
	}
	return nil
}

