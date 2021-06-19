module github.com/velann21/bloom-services/users-service

go 1.15

require (
	github.com/go-redis/redis/v8 v8.10.0
	github.com/sirupsen/logrus v1.8.1
	github.com/velann21/bloom-services/common-lib v1.0.0
)

replace github.com/velann21/bloom-services/common-lib v1.0.0 => ../go-common-lib
