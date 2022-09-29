package service

var instance Checker

type Checker interface {
	HasService(service string) (bool, error)
}

func GetChecker() Checker {
	return instance
}
