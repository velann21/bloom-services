package helpers

import (
	"bytes"
	"encoding/json"
	"time"
)

func GetDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}

func ConvertStructToBytes(data interface{})([]byte,error){
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(data)
	if err != nil{
		return nil, err
	}
	return reqBodyBytes.Bytes(), nil
}