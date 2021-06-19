package requests

import "time"

type JsonBirthDate time.Time

type User struct {
	Name string `json:"name"`
    MobileNumber string `json:"mobile_number"`
	Address Address `json:"address"`
	DOB JsonBirthDate `json:"dob"`
}

type Address struct {
	ZipCode string `json:"zip_code"`
	StreetName string `json:"street_name"`
	HouseNumber int `json:"house_number"`
}


