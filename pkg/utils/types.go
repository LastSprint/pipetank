package utils

import "time"

type Clock func() time.Time

func UTCClock() Clock {
	return func() time.Time {
		return time.Now().UTC()
	}
}
