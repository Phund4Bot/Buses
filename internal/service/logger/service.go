package logger

import (
	"fmt"
	"time"
)

const (
	logTemplate = "{'level':'%s','time':'%v','message':'%v'}\n"
	timeFormat  = "2006-01-02 15:04:05"

	infoLevel = "info"
	errorLevel = "error"
)

func Info(time time.Time, msg interface{}) {
	fmt.Printf(logTemplate, infoLevel, time.Format(timeFormat), msg)
}
func Infof(time time.Time, template string, params ...interface{}) {
	fmt.Printf(logTemplate, infoLevel, time.Format(timeFormat), fmt.Sprintf(template, params...))
}
func BuildInfoLog(time time.Time, msg interface{}) string {
	return fmt.Sprintf(logTemplate, infoLevel, time.Format(timeFormat), msg)
}
func BuildInfofLog(time time.Time, template string, params ...interface{}) string {
	return fmt.Sprintf(logTemplate, infoLevel, time.Format(timeFormat), fmt.Sprintf(template, params...))
}


func Error(time time.Time, msg interface{}) {
	fmt.Printf(logTemplate, errorLevel, time.Format(timeFormat), msg)
}
func Errorf(time time.Time, template string, params ...interface{}) {
	fmt.Printf(logTemplate, errorLevel, time.Format(timeFormat), fmt.Sprintf(template, params...))
}
func BuildErrorLog(time time.Time, msg interface{}) string {
	return fmt.Sprintf(logTemplate, errorLevel, time.Format(timeFormat), msg)
}
func BuildErrorfLog(time time.Time, template string, params ...interface{}) string {
	return fmt.Sprintf(logTemplate, errorLevel, time.Format(timeFormat), fmt.Sprintf(template, params...))
}
