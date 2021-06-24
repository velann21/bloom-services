# bloom-service
This is the Monorepo for the bloom services and currently it has the implementation for common library and user service in Go.

## Folder Structure:

  root structure:
    
    /.github       -> contains the CI/CD pipeline code)
    /go-common-lib -> contains all the common usable code across microservices)
    /users-service -> contains person microservice code)
    Makefile       -> code to do all the adhoc activities)
    
  ``
  
  .github:
  ```
   workflow/main.yaml (CI/CD file)
  ```

  go-common-lib:
  ```
   databases -> Contains the all the database client code to make it reusable
   entities  -> Contains the all request and response struct so that we can resue the struct for producer and consumer side microservices
   helpers   -> contains all the utility functions
   server    -> Mux server code to reuse the server code in multiple services
  ```
 
 users-services:
 ```
  deployments -> contains the code for user-service app deployment
  docs        -> contains the documentation related to user-service architecture things
  pkg         -> actuall server side code present
    |
    |   controller -> contains entry point for each and every API
    |   database   -> contains the databases connection init code  
    |   entities   -> contains database models and response just related to user service.
    |   repository -> contains the redis db query client
    |   routes     -> contains the mux routes
    |   service    -> contains the business logics
  main.go -> App entry point
 ```
   
# API SPEC:

## Create User API:
 ```
http://localhost:5000/users/api/v1/user

Request Body:
{
    "name": string,
    "email":string,
    "address":{
        "zip_code":string,
        "street_name":string,
        "house_number": int
    },
    "dob": {
        "year": int,
        "month": int,
        "day":int
    }
}


Sucess Response Body:
{
    "success": bool,
    "status": string,
    "data": [
        {
            "email": string,
            "message": string
        }
    ]
}


Failed Response Body:
{
    "success": false,
    "errors": [
        {
            "message": "string",
            "error_code": int
        }
    ]
}
 ```
## Get User
```
http://localhost:5000/users/api/v1/user?email=xxxx@gmail.com

Success Response Body:
{
    "success": true,
    "status": "OK",
    "data": [
        {
            "user": {
                "name": string,
                "email": string,
                "address": {
                    "zip_code": string,
                    "street_name": string,
                    "house_number": int
                },
                "dob": {
                    "year": int,
                    "month": int,
                    "day": int
                },
                "created_at": time,
                "updated_at": time
            }
        }
    ]
}

Failed Response Body:
{
    "success": false,
    "errors": [
        {
            "message": "string",
            "error_code": int
        }
    ]
}
```

## UpdateWithPessimisticLock
```
http://localhost:8086/users/api/v1/user/pessimistic
{
    "name": string,
    "email": string,
    "address":{
        "zip_code": string,
        "street_name": string,
        "house_number": int
    },
    "dob": {
        "year": int,
        "month": int,
        "day": int
    }
}

Sucess Response Body:
{
    "success": true,
    "status": "OK",
    "data": [
        {
            "email": string
        }
    ]
}

Failed Response Body:
{
    "success": false,
    "errors": [
        {
            "message": "string",
            "error_code": int
        }
    ]
}
```

## UpdateWithOptimisticLock
```
http://localhost:8086/users/api/v1/user/optimistic
{
    "name": string,
    "email": string,
    "address":{
        "zip_code": string,
        "street_name": string,
        "house_number": int
    },
    "dob": {
        "year": int,
        "month": int,
        "day": int
    }
}

Sucess Response Body:
{
    "success": true,
    "status": "OK",
    "data": [
        {
            "email": string
        }
    ]
}

Failed Response Body:
{
    "success": false,
    "errors": [
        {
            "message": "string",
            "error_code": int
        }
    ]
}
```

# PreRequisite tools to run the project 
1. GO Installed
2. Redis
3. Docker
4. Kubernetes
5. Helm

# How to run the Project:
Please use the Makefile to run the project
``````
vendor_bloom_user_service:
	To do the go mod vendor for user srv

vendor_bloom_common_lib:
	To do the go mod vendor for go common lib srv

build_bloom_user_service:vendor_bloom_user_service
	To build docker image for user service

push_bloom_user_service:build_bloom_user_service
	To push the doker image to docker hub

deploy_bloom_user_service:
	To deploy apps to kubernetes

unit_test_bloom_user_service:
	To run the user service unit test

integration_test_bloom_user_service:
	To run the user service integrration test

```````

# Should be improved/Implemented:
```
1. Should Prometheus metrics
2. Should implement the Bonus point writting the expired user to filesystem
3. Add 100% test coverage
4. Distributed lock for clusterd redis
```

## Assumption:
1. I assumed that the there will be single instance of redis so I haven't used distributed lock(Redlock)

## Notes:
1. I have completed the assignment in almost 8hrs, I could improve lot of things in this assignment.
2. Due to time constraints I have just implemented the unit test and integration test only for controller and service layer

## How Pessimistic Lock Implemented:
```
1. 
```