package log

import (
	"testing"
)

func Benchmark_JsonLogSimpleType(b *testing.B) {
	Logger := NewLogger().SetWriters(NewFileWriter().Filename(logfilename).RotateHour())
	defer Logger.Close()

	for i := 0; i < b.N; i++ { //use b.N for looping
		Logger.Info(i)
	}
}

func Benchmark_JsonLogStruct(b *testing.B) {
	Logger := NewLogger().SetWriters(NewFileWriter().Filename(logfilename).RotateHour())
	defer Logger.Close()

	p1 := LoggerTestPlayer{
		1,
		"p1",
		[]LoggerTestItem{
			LoggerTestItem{1001, 1, 10},
			LoggerTestItem{1002, 23, 1},
		},
	}

	for i := 0; i < b.N; i++ { //use b.N for looping
		Logger.Info(p1)
	}
}

func Benchmark_LogSimpleType(b *testing.B) {
	Logger := NewLogger().SetWriters(NewFileWriter().Filename(logfilename).RotateHour()).UseFlatFormat()
	defer Logger.Close()

	for i := 0; i < b.N; i++ { //use b.N for looping
		Logger.Info(i)
	}
}

func Benchmark_LogStruct(b *testing.B) {
	Logger := NewLogger().SetWriters(NewFileWriter().Filename(logfilename).RotateHour()).UseFlatFormat()
	defer Logger.Close()

	p1 := LoggerTestPlayer{
		1,
		"p1",
		[]LoggerTestItem{
			LoggerTestItem{1001, 1, 10},
			LoggerTestItem{1002, 23, 1},
		},
	}

	for i := 0; i < b.N; i++ { //use b.N for looping
		Logger.Info(p1)
	}
}
