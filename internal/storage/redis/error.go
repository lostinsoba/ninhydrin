package redis

import (
	"database/sql"

	"github.com/lib/pq"
)

const (
	errCodeAlreadyExist = "23505"
)

func isAlreadyExist(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == errCodeAlreadyExist
	}
	return false
}

func isNoRows(err error) bool {
	return err == sql.ErrNoRows
}
