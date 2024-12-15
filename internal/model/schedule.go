package model

type WorkDay struct {
	Day       int    `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Schedule struct {
	WorkTime  []WorkDay `json:"work_time"`
	BreakTime []WorkDay `json:"break_time"`
}
