package model

type ErrNotFound struct{}

func (ErrNotFound) Error() string {
	return "not found"
}

type ErrAlreadyExist struct{}

func (ErrAlreadyExist) Error() string {
	return "already exist"
}
