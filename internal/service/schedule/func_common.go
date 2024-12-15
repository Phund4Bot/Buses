package schedule

import (
	"fmt"
	"math/rand"
	"time"

	"kurs/internal/model"
	"kurs/internal/utils"
)

const (
	type1SchedulePath = "./data/schedule_type1.json"
	type2SchedulePath = "./data/schedule_type2.json"
)

func (s *service) FindNextDriver(currentTime time.Time) *model.Driver {
	for id, driver := range s.drivers {
		if driver.IsAvailable(currentTime) {
			delete(s.drivers, id)
			return driver
		}
	}
	return nil
}

func isPeakHour(currentTime time.Time) bool {
	hour := currentTime.Hour()
	return (7 <= hour && hour <= 9) || (17 <= hour && hour <= 19)
}
func isNightTime(currentTime time.Time) bool {
	hour := currentTime.Hour()
	return hour <= 6 || hour >= 22
}

func calculatePeopleCount(currentTime time.Time) int {
	switch {
	case isNightTime(currentTime) && rand.Float32() < 0.2:
		return rand.Intn(2) + 1
	case isPeakHour(currentTime):
		return rand.Intn(3) + 3
	case rand.Float32() < 0.3:
		return rand.Intn(3) + 1
	default:
		return 0
	}
}

func selectDriverSchedule(scheduleType int) (model.Schedule, error) {
	switch scheduleType {
	case 1:
		schedule, err := utils.ParseFromFile[model.Schedule](type1SchedulePath)
		if err != nil {
			return schedule, fmt.Errorf("cannot parse schedule: %v", err)
		}
		return schedule, nil
	default:
		schedule, err := utils.ParseFromFile[model.Schedule](type2SchedulePath)
		if err != nil {
			return schedule, fmt.Errorf("cannot parse schedule: %v", err)
		}
		return schedule, nil
	}
}

func selectRandomIndex(probabilities []float64) int {
	r := rand.Float64()
	sum := 0.0
	for i, p := range probabilities {
		sum += p
		if r <= sum {
			return i
		}
	}
	return len(probabilities) - 1
}
