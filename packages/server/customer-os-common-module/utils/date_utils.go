package utils

import "time"

func CurrentMonth() string {
	currentTime := time.Now()
	return currentTime.Format("2006-01")
}
