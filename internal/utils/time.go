package utils

import "time"

func GetNextMonday(startDate time.Time) time.Time {
	if startDate.IsZero() {
		startDate = time.Now()
	}
	currentWeekday := startDate.Weekday()

	daysUntilNextMonday := (8 - int(currentWeekday)) % 7
	if daysUntilNextMonday == 0 {
		daysUntilNextMonday = 7
	}
	
	return startDate.AddDate(0, 0, daysUntilNextMonday)
}
