package tools

import (
	"time"
)

var (
	formatedTime string
)

func CurrentTime() string {
	t := time.Now()
	formatedTime = t.Format(time.RFC1123)
	return formatedTime
}
