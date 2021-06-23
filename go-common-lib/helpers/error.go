package helpers

import "errors"

var (
	NoResultFound        = errors.New("NoResultsFound")
	InvalidRequest       = errors.New("InvalidRequest")
	SomethingWrong       = errors.New("SomethingWentWrong")
	UserAlreadyExist     = errors.New("UserExists")
	RedisNil             = "redis: nil"
	RedisNilError        = errors.New("redis: nil")
	TryLater             = errors.New("AlreadyOneUpdateInProgress/TryLater")
	UpdateNotPossibleNow = errors.New("UpdateNotPossible")
)
