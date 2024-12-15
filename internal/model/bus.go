package model

import (
	"time"
)

type Bus struct {
	BusID int `json:"bus_id"`

	Capacity          int     `json:"capacity"`
	IsActive          bool    `json:"is_active"`
	Driver            *Driver `json:"driver"`
	PassengersCounter int     `json:"passengers_counter"`

	StartBusStop int `json:"start_bus_stop"`
	StartTime    time.Time
}

func (b *Bus) SetDriver(driver *Driver) {
	b.Driver = driver
}
func (b *Bus) SetActive(new bool) {
	b.IsActive = new
}

func (b *Bus) CheckDriverAvailability(currentTime time.Time) {
	if b.Driver == nil {
		b.IsActive = false
		return
	}

	if !b.Driver.IsAvailable(currentTime) {
		b.IsActive = false
	} else {
		b.IsActive = true
	}
}

type BusStop struct {
	Name          string            `json:"name"`
	Duration      int               `json:"duration"`
	TimeToNext    int               `json:"time_to_next"`
	PeopleWaiting map[int64][]Human `json:"people_waiting"`
}

func (s *BusStop) GetPeopleCount() int {
	count := 0

	for _, value := range s.PeopleWaiting {
		count += len(value)
	}

	return count
}

func (s *BusStop) UnloadPeople(count int) {
	for key, value := range s.PeopleWaiting {
		if count == 0 {
			break
		}

		l := len(value)
		if l >= count {
			s.PeopleWaiting[key] = value[count:]
			count = 0
		} else {
			s.PeopleWaiting[key] = make([]Human, len(value))
			count -= l
		}
	}
}
