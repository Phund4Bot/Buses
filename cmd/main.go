package main

import (
	"time"

	"kurs/internal/service/logger"
	"kurs/internal/service/schedule"
)

func main() {
	driverTypesSlice := []int{1, 1, 1, 1, 1, 1, 2}
	scheduleService, err := schedule.Initialize(driverTypesSlice)
	if err != nil {
		logger.Errorf(time.Now(), "Cannot init schedule service: %v", err)
		return
	}
	pass := scheduleService.RunSimulation()
	logger.Infof(time.Now(), "Test simulation result: %v", pass)

	brutforce := schedule.NewBrutforce()
	brutforce.RunBrutforce(len(driverTypesSlice))

	genetic := schedule.NewGenetic(10, 5, 0.1)
	solution := genetic.RunGenetic()
	logger.Infof(time.Now(), "Best solution: %v", solution)
}
