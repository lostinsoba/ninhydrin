package model

type ErrNotFound struct{}

func (ErrNotFound) Error() string {
	return "not found"
}
