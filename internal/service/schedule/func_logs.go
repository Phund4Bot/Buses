package schedule

import (
	"fmt"
	"sort"
	"text/tabwriter"

	"kurs/internal/model"
	"kurs/internal/utils"
)

const (
	tablePath = "./logs/table.log"
	logsPath  = "./logs/logs.log"
)

func (s *service) appendLog(log string) {
	s.scheduleLogs = append(s.scheduleLogs, log)
}

func (s *service) saveLogs(logs []string) error {
	file, err := utils.CreateFileConnection(logsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, log := range logs {
		_, err = file.WriteString(log)
		if err != nil {
			return fmt.Errorf("cannot write string: %v", err)
		}
	}

	return nil
}

func (s *service) saveTable(rows []model.Row) error {
	file, err := utils.CreateFileConnection(tablePath)
	if err != nil {
		return err
	}
	defer file.Close()

	BusIDMap := make(map[int][]model.Row, len(s.buses))

	for _, row := range rows {
		if _, ok := BusIDMap[row.BusID]; !ok {
			BusIDMap[row.BusID] = make([]model.Row, 0, len(rows)/len(s.buses))
		}
		BusIDMap[row.BusID] = append(BusIDMap[row.BusID], row)
	}

	writer := tabwriter.NewWriter(file, 0, 3, 2, ' ', tabwriter.Debug)

	for busID, value := range BusIDMap {
		fmt.Fprintf(writer, "BusID: %d\n", busID)
		fmt.Fprintln(writer, "Time\tStopName\tDriverID\tPassengersCount\tBusLoad")

		sort.Slice(value, func(i, j int) bool {
			return value[i].Time < value[j].Time
		})

		for _, row := range value {
			_, err := fmt.Fprintf(writer, "%s\t%s\t%d\t%d\t%d%%\n", row.Time, row.StopName, row.DriverID, row.PassengersCount, row.BusLoad)
			if err != nil {
				return fmt.Errorf("cannot write logs to file: %v", err)
			}
		}

		fmt.Fprintln(writer, "------------------------------------------------------------------------")
	}

	return writer.Flush()
}
