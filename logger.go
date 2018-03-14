package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	DebugFilter = 1
	TraceFilter = 2
	InfoFilter  = 4
	WarnFilter  = 8
	ErrorFilter = 16
	FatalFilter = 32
)

// var levStrings = []string{"DEBUG", "TRACE", "INFO", "WARN", "ERROR", "FATAL"}

type formatter interface {
	Format(level, file, callstack string, msg interface{}) string
}

//Logger logger struct
type Logger struct {
	wait    sync.WaitGroup
	logchan chan string
	closed  bool
	writers []io.Writer
	filters int
	formatter
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//NewLogger create a new logger
func NewLogger() *Logger {
	logger := &Logger{
		logchan:   make(chan string, 30),
		closed:    false,
		formatter: &FlatFormatter{},
		filters:   0,
	}

	logger.writers = append(logger.writers, os.Stdout)

	logger.wait.Add(1)
	go logger.routine()
	return logger
}

func (l *Logger) DebugOn() *Logger {
	l.filters ^= DebugFilter
	return l
}

func (l *Logger) DebugOff() *Logger {
	l.filters |= DebugFilter
	return l
}

func (l *Logger) TraceOn() *Logger {
	l.filters ^= TraceFilter
	return l
}

func (l *Logger) TraceOff() *Logger {
	l.filters |= TraceFilter
	return l
}

func (l *Logger) InfoOn() *Logger {
	l.filters ^= InfoFilter
	return l
}

func (l *Logger) InfoOff() *Logger {
	l.filters |= InfoFilter
	return l
}

func (l *Logger) WarnOn() *Logger {
	l.filters ^= WarnFilter
	return l
}

func (l *Logger) WarnOff() *Logger {
	l.filters |= WarnFilter
	return l
}

func (l *Logger) ErrorOn() *Logger {
	l.filters ^= ErrorFilter
	return l
}

func (l *Logger) ErrorOff() *Logger {
	l.filters |= ErrorFilter
	return l
}

func (l *Logger) FatalOn() *Logger {
	l.filters ^= FatalFilter
	return l
}

func (l *Logger) FatalOff() *Logger {
	l.filters |= FatalFilter
	return l
}

//AppendWriters append custom writers
func (l *Logger) AppendWriters(writers ...io.Writer) *Logger {
	l.writers = append(l.writers, writers...)
	return l
}

//SetWriters set custom writers
func (l *Logger) SetWriters(writers ...io.Writer) *Logger {
	l.writers = writers
	return l
}

//UseFlatFormat use flat format
func (l *Logger) UseFlatFormat() *Logger {
	l.formatter = &FlatFormatter{}
	return l
}

//UseJSONFormat use json format
func (l *Logger) UseJSONFormat(indent bool) *Logger {
	l.formatter = &JSONFormatter{
		indent: indent,
	}
	return l
}

//UseFormatter set custom formatter
func (l *Logger) UseFormatter(f formatter) *Logger {
	l.formatter = f
	return l
}

//Close close logger
func (l *Logger) Close() {
	l.closed = true
	close(l.logchan)
	l.wait.Wait()
}

//Debug Write log use debug level
func (l *Logger) Debug(v interface{}, args ...interface{}) {
	if l.filters&DebugFilter > 0 {
		return
	}
	l.writeLog("DEBUG", l.getFileAndLine(), "", v, args...)
}

//Trace Write log use trace level
func (l *Logger) Trace(v interface{}, args ...interface{}) {
	if l.filters&TraceFilter > 0 {
		return
	}
	l.writeLog("TRACE", l.getFileAndLine(), "", v, args...)
}

//Info Write log use info level
func (l *Logger) Info(v interface{}, args ...interface{}) {
	if l.filters&InfoFilter > 0 {
		return
	}
	l.writeLog("INFO ", l.getFileAndLine(), "", v, args...)
}

//Warn Write log use warn level
func (l *Logger) Warn(v interface{}, args ...interface{}) {
	if l.filters&WarnFilter > 0 {
		return
	}
	l.writeLog("WARN ", l.getFileAndLine(), "", v, args...)
}

//Warn Write log use error level
func (l *Logger) Error(v interface{}, args ...interface{}) {
	if l.filters&ErrorFilter > 0 {
		return
	}
	l.writeLog("ERROR", l.getFileAndLine(), string(debug.Stack()), v, args...)
}

//Fatal Write log use fatal level
func (l *Logger) Fatal(v interface{}, args ...interface{}) {
	if l.filters&FatalFilter > 0 {
		return
	}
	l.writeLog("FATAL", l.getFileAndLine(), string(debug.Stack()), v, args...)
}

func (l *Logger) routine() {
	defer l.wait.Done()

	exit := false

	for !exit {
		select {
		case str, ok := <-l.logchan:
			if !ok {
				exit = true
				break
			}
			for _, writer := range l.writers {
				writer.Write([]byte(str))
				writer.Write([]byte("\n"))
			}
		}
	}
}

func (l *Logger) getFileAndLine() string {
	var ok bool
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	return fmt.Sprintf("%s:%d", path.Base(file), line)
}

func (l *Logger) writeLog(levstr string, filestr string, stackstr string, v interface{}, args ...interface{}) {
	if l.closed {
		return
	}

	var msg interface{}
	s, ok := v.(string)
	if ok {
		msg = fmt.Sprintf(s, args...)
	} else {
		msg = v
	}
	l.logchan <- l.formatter.Format(levstr, filestr, stackstr, msg)
}
