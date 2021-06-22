package helpers

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
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

func DetectAppMode(args []string){
	for _, arg := range args{
		if arg == "Dev"{
			lvl := "debug"
			logrus.Info("App running in Dev Mode")
			SetEnvUsersDevelopmentMode()
			ll, err := logrus.ParseLevel(lvl)
			if err != nil {
				ll = logrus.DebugLevel
			}
			logrus.SetLevel(ll)
			break
		}else if arg == "Prod"{
			logrus.Info("App running in Prod Mode")
			if GetEnv(REDIS) == ""{
				logrus.Info("Redis conn string cannot be empty")
				os.Exit(1)
			}
			break
		}
	}
}
