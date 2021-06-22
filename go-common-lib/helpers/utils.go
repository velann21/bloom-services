package helpers

import (
	"bytes"
	"encoding/json"
	"os"
	"time"
)

func GetDOB(year, month, day int) time.Time {
	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return dob
}

func ConvertStructToBytes(data interface{}) ([]byte, error) {
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(data)
	if err != nil {
		return nil, err
	}
	return reqBodyBytes.Bytes(), nil
}

const (
	REDIS = "REDIS_CONN"
)

func GetEnv(key string) string {
	switch key {
	case REDIS:
		return os.Getenv(REDIS)
	}
	return ""
}

func SetEnvUsersDevelopmentMode() {
	_ = os.Setenv("REDIS_CONN", "127.0.0.1:6379")
}
