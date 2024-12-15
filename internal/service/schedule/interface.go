package schedule

import "kurs/internal/model"

type Service interface {
	AddBus(bus *model.Bus, startDelay, startStop int)
	AddBusStop(busStop *model.BusStop)
	AddDriver(driver *model.Driver)

	RunSimulation() int
}

type Genetic interface {
	RunGenetic() []int
}

type Brutforce interface {
	RunBrutforce(driversCount int)
}
