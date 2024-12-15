package model

import (
	"time"
)

type Driver struct {
	DriverID  int       `json:"driver_id"`
	IsActive  bool      `json:"is_active"`
	WorkTime  []WorkDay `json:"work_time"`
	BreakTime []WorkDay `json:"break_time"`
}

func (d *Driver) SetActive(new bool) {
	d.IsActive = new
}
func (d *Driver) SetWorkTime(workTime []WorkDay) {
	d.WorkTime = workTime
}
func (d *Driver) SetBreakTime(breakTime []WorkDay) {
	d.BreakTime = breakTime
}

func (d *Driver) IsAvailable(currentTime time.Time) bool {
	currentWeekday := int(currentTime.Weekday())
	currentDate := currentTime.Format("2006-01-02")

	for _, work := range d.WorkTime {
		if work.Day == currentWeekday {
			workStart, _ := time.Parse("2006-01-02 15:04", currentDate+" "+work.StartTime)
			workEnd, _ := time.Parse("2006-01-02 15:04", currentDate+" "+work.EndTime)

			if currentTime.After(workStart) && currentTime.Before(workEnd) {
				for _, breakTime := range d.BreakTime {
					if breakTime.Day == currentWeekday {
						breakStart, _ := time.Parse("2006-01-02 15:04", currentDate+" "+breakTime.StartTime)
						breakEnd, _ := time.Parse("2006-01-02 15:04", currentDate+" "+breakTime.EndTime)

						if currentTime.After(breakStart) && currentTime.Before(breakEnd) {
							return false
						}
					}
				}
				return true
			}
		}
	}
	return false
}
