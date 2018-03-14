package log

import (
	"fmt"
	"os"
	"path"
	"time"
)

const (
	RotMin  = 1
	RotHour = 2
	RotDay  = 3
)

type FileWriter struct {
	filename   string
	logfile    *os.File
	logdatestr string
	rotate     int
}

func NewFileWriter() *FileWriter {
	w := &FileWriter{
		rotate:   RotDay,
		filename: path.Join(getCurrentDirectory(), "logs/default.log"),
	}
	return w
}

func (w *FileWriter) Filename(filename string) *FileWriter {
	if filename == "" {
		return w
	}
	w.filename = filename
	return w
}

func (w *FileWriter) RotateMin() *FileWriter {
	w.rotate = RotMin
	return w
}

func (w *FileWriter) RotateHour() *FileWriter {
	w.rotate = RotHour
	return w
}

func (w *FileWriter) RotateDay() *FileWriter {
	w.rotate = RotDay
	return w
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	err = w.updateLogFile()
	if err != nil {
		return 0, err
	}

	return w.logfile.Write(p)
}

func (w *FileWriter) getRotateDateStr() (datestr string) {
	switch w.rotate {
	case RotMin:
		datestr = time.Now().Format("2006_01_02_15_04")
	case RotHour:
		datestr = time.Now().Format("2006_01_02_15")
	default:
		datestr = time.Now().Format("2006_01_02")
	}

	return datestr
}

func (w *FileWriter) openLogFile() (logfile *os.File, err error) {
	dirstr, filestr := path.Split(w.filename)

	filestr = fmt.Sprintf("%s_%s", w.getRotateDateStr(), filestr)
	fullpath := path.Join(dirstr, filestr)

	err = os.MkdirAll(dirstr, 0666)
	if err != nil {
		return logfile, err
	}

	logfile, err = os.OpenFile(fullpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return logfile, err
	}

	return logfile, err
}

func (w *FileWriter) updateLogFile() (err error) {
	datestr := w.getRotateDateStr()
	if w.logdatestr != datestr {
		w.logfile.Close()
		w.logfile = nil
	}

	if w.logfile == nil {
		w.logfile, err = w.openLogFile()
		if err != nil {
			return err
		}
		w.logdatestr = datestr
	}

	return nil
}
