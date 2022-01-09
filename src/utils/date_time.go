package utils

import "time"

var ISO8601 = "2006-01-02T15:04:05.000Z"

func GetCurrentTimeStamp() time.Time {
	return time.Now().UTC()
}

func GetCurrentISOTime() string {
	return time.Now().UTC().Format(ISO8601)
}
