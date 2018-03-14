package log

import (
	"fmt"
	"time"
)

//FlatFormatter format log with flat text
type FlatFormatter struct {
}

//Format convert log to flat text
func (f *FlatFormatter) Format(level, file, callstack string, msg interface{}) string {
	date := time.Now().Format("2006-01-02 15:04:05")
	if len(callstack) > 0 {
		return fmt.Sprintf("%s | %s | %s | %+v\n%s", date, level, file, msg, callstack)
	} else if len(file) > 0 {
		return fmt.Sprintf("%s | %s | %s | %+v", date, level, file, msg)
	} else {
		return fmt.Sprintf("%s | %s | %s | %+v", date, level, file, msg)
	}
}
