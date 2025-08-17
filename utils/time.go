package utils

import "time"

func Delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
