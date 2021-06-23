vendor_bloom_user_service:
	cd ./users-service && go mod vendor

vendor_bloom_common_lib:
	cd ./go-common-lib && go mod vendor

build_bloom_user_service:vendor_bloom_user_service
	cd ./users-service && docker build . -t bloom_user_service:$(tag) && docker tag  bloom_user_service:$(tag) singaravelan21/bloom_user_service:$(tag)

push_bloom_user_service:build_bloom_user_service
	docker push singaravelan21/bloom_user_service:$(tag)

build_custom_redis:
	cd ./users-service/deployments/docker && docker build -f DockerfileRedis . -t singaravelan21/redis

push_custom_redis:build_custom_redis
	docker push singaravelan21/redis

deploy_bloom_user_service:
	cd ./users-service/deployments/helm_charts && helm upgrade --install redis ./redis --atomic --timeout 5m0s -n bloomreach && helm upgrade --install bloom-user-services ./user-service --set image.tag=$(tag) --atomic --timeout 5m0s -n bloomreach

unit_test_bloom_user_service:
	cd ./users-service && go test -v ./... -v -short

integration_test_bloom_user_service:
	cd ./users-service && go test -v ./... -v -run Integration


