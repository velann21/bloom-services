module github.com/velann21/bloom-services/users-service

go 1.15

require (
	github.com/go-redis/redis/v8 v8.10.0
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.8.1
	github.com/velann21/bloom-services/common-lib v1.3.13 // indirect
	github.com/velann21/bloom-users-service v0.0.0-20210619130712-7486cfdef93f
)

replace github.com/velann21/bloom-services/common-lib v1.3.13 => ../go-common-lib
