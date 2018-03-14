package log

import (
	"testing"
)

type LoggerTestItem struct {
	ItemID int64
	Type   int
	Value  int
}

type LoggerTestPlayer struct {
	PlayerID   int64
	PlayerName string
	Items      []LoggerTestItem
}

var logfilename = "F:/MyProjects/code/src/moton/go-log/logs/log.txt"

func Test_Logger(t *testing.T) {
	Logger := NewLogger().AppendWriters(NewFileWriter().Filename(logfilename).RotateHour())
	defer Logger.Close()

	p1 := LoggerTestPlayer{
		1,
		"p1",
		[]LoggerTestItem{
			LoggerTestItem{1001, 1, 10},
			LoggerTestItem{1002, 23, 1},
		},
	}

	p2 := LoggerTestPlayer{
		2,
		"p2",
		[]LoggerTestItem{
			LoggerTestItem{1005, 50, 1000},
		},
	}

	Logger.UseJSONFormat(true)
	Logger.Debug("Hello World! %s!", "My Logger!")
	Logger.Trace("Look here")
	Logger.Info(p1)
	Logger.Warn(15999)
	Logger.Error(0.123129321)
	Logger.Fatal(p2)

	Logger.UseFlatFormat()
	Logger.Debug("Hello World! %s!", "My Logger!")
	Logger.Trace("Look here")
	Logger.Info(p1)
	Logger.Warn(15999)
	Logger.Error(0.123129321)
	Logger.Fatal(p2)

	Logger.DebugOff().FatalOff()
	Logger.Debug("Hello World! %s!", "My Logger!")
	Logger.Trace("Look here")
	Logger.Info(p1)
	Logger.Warn(15999)
	Logger.Error(0.123129321)
	Logger.Fatal(p2)

	Logger.DebugOn().FatalOn()
	Logger.Debug("Hello World! %s!", "My Logger!")
	Logger.Trace("Look here")
	Logger.Info(p1)
	Logger.Warn(15999)
	Logger.Error(0.123129321)
	Logger.Fatal(p2)
}
