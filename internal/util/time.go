package util

import "time"

func UnixEpoch() int64 {
	return time.Now().UTC().Unix()
}
