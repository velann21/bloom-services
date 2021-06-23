package helpers

import "errors"

var (
	NoResultFound        = errors.New("NoResultsFound")
	InvalidRequest       = errors.New("InvalidRequest")
	SomethingWrong       = errors.New("SomethingWentWrong")
	UserAlreadyExist     = errors.New("UserExists")
	RedisNil             = "redis: nil"
	RedisNilError        = errors.New("redis: nil")
	UpdateNotPossibleNow = errors.New("UpdateNotPossible")
)
