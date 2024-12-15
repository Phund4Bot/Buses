package schedule

import (
	"kurs/internal/service/logger"
	"kurs/internal/utils"
	"time"
)

type brutforce struct {
}

func NewBrutforce() Brutforce {
	return &brutforce{}
}

func generateAllCombinations(current []int, length int, result *[][]int) {
	if len(current) == length {
		combination := make([]int, length)
		copy(combination, current)
		*result = append(*result, combination)
		return
	}

	for _, val := range []int{1, 2} {
		current = append(current, val)
		generateAllCombinations(current, length, result)
		current = current[:len(current)-1]
	}
}

func (b *brutforce) RunBrutforce(driversCount int) {
	logger.Info(time.Now(), "Starting brutforce method")
	defer utils.LogElapsed("Finished brutforce method")()

	driverVariants := make([][]int, 0, 1<<8)
	generateAllCombinations([]int{}, driversCount, &driverVariants)

	var (
		bestSolution []int
		bestResult   int
	)

	for _, drivers := range driverVariants {
		schedule, err := Initialize(drivers)
		if err != nil {
			logger.Errorf(time.Now(), "Cannot inititalize schedule: %v", err)
			return 
		}
		
		result := schedule.RunSimulation()
		if result > bestResult {
			bestResult = result
			bestSolution = drivers
		}
	}

	logger.Infof(time.Now(), "Brutforce statistic: best result: %v, best solution: %v", bestResult, bestSolution)
}
