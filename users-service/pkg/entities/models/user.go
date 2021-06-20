package models

import (
	"bytes"
	"encoding/gob"
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

func (user *User) PopulateUser(result []byte) error {
	ioReader := bytes.NewReader(result)
	err := gob.NewDecoder(ioReader).Decode(user)
	if err != nil{
		return err
	}
	return nil
}

