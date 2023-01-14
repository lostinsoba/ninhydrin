package redis

import r "github.com/go-redis/redis/v9"

func isNil(err error) bool {
	return err == r.Nil
}
