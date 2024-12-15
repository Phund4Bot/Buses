package utils

import (
	"time"

	"kurs/internal/service/logger"
)

func LogElapsed(msg string) func() {
	started := time.Now()
	return func() {
		logger.Infof(time.Now(), "%v %v elapsed (sec)", msg, time.Since(started).Seconds())
	}
}
