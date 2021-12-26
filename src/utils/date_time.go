package utils

import "time"

func GetCurrentTimeStamp() time.Time {
	return time.Now().UTC()
}
