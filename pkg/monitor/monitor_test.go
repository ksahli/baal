package monitor_test

import (
	"time"
)

var timestamp = time.Now().UTC()

func stamper() time.Time {
	return timestamp
}
