package schedule

import (
	"fmt"
	"math/rand"
	"time"

	"kurs/internal/model"
	"kurs/internal/service/logger"
	"kurs/internal/utils"
)

const (
	busStopsDataPath = "./data/bus_stops.json"
	driversDataPath  = "./data/drivers.json"
	busesDataPath    = "./data/buses.json"
)

type service struct {
	StartTime time.Time
	EndTime   time.Time
	TimeStep  time.Duration

	busStops []*model.BusStop
	buses    []*model.Bus
	drivers  map[int]*model.Driver

	passengersCounter     int
	leftPassengersCounter int

	scheduleRows []model.Row
	scheduleLogs []string
}

func New(startTime, endTime time.Time, timeStep time.Duration) Service {
	return &service{
		StartTime: startTime,
		EndTime:   endTime,
		TimeStep:  timeStep,

		busStops: make([]*model.BusStop, 0),
		buses:    make([]*model.Bus, 0),
		drivers:  make(map[int]*model.Driver),

		scheduleRows: make([]model.Row, 0, (endTime.Unix()-startTime.Unix())/int64(timeStep.Seconds())),
		scheduleLogs: make([]string, 0, (endTime.Unix()-startTime.Unix())/int64(timeStep.Seconds())),
	}
}

func Initialize(driverTypes []int) (Service, error) {
	startDate := utils.GetNextMonday(time.Now())
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 9, 0, 0, 0, startDate.Location())
	endTime := startTime.AddDate(0, 0, 7)

	schedule := New(startTime, endTime, time.Minute)

	busStops, err := utils.ParseFromFile[[]model.BusStop](busStopsDataPath)
	if err != nil {
		return nil, fmt.Errorf("cannot parse bus stops: %v", err)
	}
	for _, busStop := range busStops {
		busStop.PeopleWaiting = make(map[int64][]model.Human, 10)
		schedule.AddBusStop(&busStop)
	}

	drivers, err := utils.ParseFromFile[[]model.Driver](driversDataPath)
	if err != nil {
		return nil, fmt.Errorf("cannot parse drivers: %v", err)
	}
	for i := 0; i < len(drivers); i++ {
		driverSchedule, err := selectDriverSchedule(driverTypes[i])
		if err != nil {
			return nil, err
		}
		drivers[i].SetWorkTime(driverSchedule.WorkTime)
		drivers[i].SetBreakTime(driverSchedule.BreakTime)
		schedule.AddDriver(&drivers[i])
	}

	buses, err := utils.ParseFromFile[[]model.Bus](busesDataPath)
	if err != nil {
		return nil, fmt.Errorf("cannot parse buses: %v", err)
	}
	for i, bus := range buses {
		schedule.AddBus(&bus, i*5, i%len(busStops))
	}

	return schedule, nil
}

func (s *service) AddBus(bus *model.Bus, startDelay, startBusStop int) {
	startTime := s.StartTime.Add(time.Duration(startDelay) * time.Minute)
	bus.StartBusStop = startBusStop
	bus.StartTime = startTime
	s.buses = append(s.buses, bus)
}
func (s *service) AddBusStop(busStop *model.BusStop) {
	s.busStops = append(s.busStops, busStop)
}
func (s *service) AddDriver(driver *model.Driver) {
	s.drivers[driver.DriverID] = driver
}

func (s *service) RunSimulation() int {
	logger.Info(time.Now(), "Starting simulation")
	defer utils.LogElapsed("Finished simulation")()

	s.appendLog(logger.BuildInfoLog(s.StartTime, "Starting bus scheduling"))

	busPositions := make(map[int]int)
	busDepartureTimes := make(map[int]time.Time)

	for _, bus := range s.buses {
		busPositions[bus.BusID] = bus.StartBusStop
		busDepartureTimes[bus.BusID] = bus.StartTime
	}

	for currentTime := s.StartTime; currentTime.Before(s.EndTime); currentTime = currentTime.Add(s.TimeStep) {
		for _, busStop := range s.busStops {

			peopleCount := calculatePeopleCount(currentTime)
			people := model.GeneratePeople(peopleCount)
			for _, human := range people {
				timestamp := currentTime.Add(human.WaitingTime).Unix()
				if _, ok := busStop.PeopleWaiting[timestamp]; !ok {
					busStop.PeopleWaiting[timestamp] = make([]model.Human, 0, len(people))
				}
				busStop.PeopleWaiting[timestamp] = append(busStop.PeopleWaiting[timestamp], human)
			}

			for key, value := range busStop.PeopleWaiting {
				if time.Unix(key, 0).Compare(currentTime) == 0 {
					s.leftPassengersCounter += len(value)
					delete(busStop.PeopleWaiting, key)
				}
			}
		}

		for _, bus := range s.buses {
			if !bus.IsActive {
				nextDriver := s.FindNextDriver(currentTime)
				if nextDriver != nil {
					bus.Driver = nextDriver
					bus.SetActive(true)
					s.appendLog(logger.BuildInfofLog(currentTime, "Bus %d is active with driver %d", bus.BusID, bus.Driver.DriverID))
				}
				continue
			}

			if !bus.Driver.IsAvailable(currentTime) {
				s.drivers[bus.Driver.DriverID] = bus.Driver
				s.appendLog(logger.BuildInfofLog(currentTime, "Bus %d has finished moving (driver %d has finished workday)", bus.BusID, bus.Driver.DriverID))
				bus.Driver = nil
			}

			bus.CheckDriverAvailability(currentTime)
			if !bus.IsActive {
				continue
			}

			currentPos := busPositions[bus.BusID]
			if currentTime.After(busDepartureTimes[bus.BusID]) {
				busStop := s.busStops[currentPos]
				freeSpace := bus.Capacity - bus.PassengersCounter
				passengersOnBoard := min(freeSpace, busStop.GetPeopleCount())

				bus.PassengersCounter += passengersOnBoard
				busStop.UnloadPeople(passengersOnBoard)
				s.passengersCounter += passengersOnBoard

				passengersOutBoard := rand.Intn(bus.PassengersCounter/3 + 1)
				bus.PassengersCounter -= passengersOutBoard
				s.appendLog(logger.BuildInfofLog(currentTime, "Bus %d has coming to %s. Passengers: +%d, -%d. Load: %d/%d.",
					bus.BusID, busStop.Name, passengersOnBoard, passengersOutBoard, bus.PassengersCounter, bus.Capacity))

				busDepartureTimes[bus.BusID] = currentTime.Add(time.Duration(busStop.Duration) * time.Minute)
				busPositions[bus.BusID] = (currentPos + 1) % len(s.busStops)
				busDepartureTimes[bus.BusID] = busDepartureTimes[bus.BusID].Add(time.Duration(busStop.TimeToNext) * time.Minute)

				s.scheduleRows = append(s.scheduleRows, model.Row{
					StopName:        busStop.Name,
					Time:            currentTime.Format("2006-01-02 15:04:05"),
					BusID:           bus.BusID,
					DriverID:        bus.Driver.DriverID,
					PassengersCount: bus.PassengersCounter,
					BusLoad:         int((float64(bus.PassengersCounter) / float64(bus.Capacity)) * 100),
				})
			}
		}
	}

	s.appendLog(logger.BuildInfofLog(s.EndTime, "Ending simulation. Transported: %d. Missed: %d.", s.passengersCounter, s.leftPassengersCounter))

	if err := s.saveLogs(s.scheduleLogs); err != nil {
		logger.Errorf(time.Now(), "Cannot save logs to file: %v", err)
	}

	if err := s.saveTable(s.scheduleRows); err != nil {
		logger.Errorf(time.Now(), "Cannot save table to file: %v", err)
	}

	return s.passengersCounter
}
