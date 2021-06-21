vendor_bloom_user_service:
	cd ./users-service && go mod vendor

vendor_bloom_common_lib:
	cd ./go-common-lib && go mod vendor

build_bloom_user_service:vendor_bloom_user_service
	cd ./users-service && docker build . -t bloom_user_service:$(tag) && docker tag  bloom_user_service:$(tag) singaravelan21/bloom_user_service:$(tag)

push_bloom_user_service:build_bloom_user_service
	docker push singaravelan21/bloom_user_service:$(tag)
