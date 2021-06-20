package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Address      Address   `json:"address"`
	DOB          DOB       `json:"dob"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

func (user *User) PopulateUser(body []byte) error {
	fmt.Println(string(body))
	err := json.Unmarshal(body, user)
	if err != nil {
		return err
	}
	return nil
}

