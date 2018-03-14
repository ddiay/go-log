package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

//JSONFormatter format log with json
type JSONFormatter struct {
	Date      string      `json:"date"`
	Lev       string      `json:"lev"`
	Msg       interface{} `json:"msg"`
	File      string      `json:"file,omitempty"`
	Callstack []string    `json:"callstack,omitempty"`
	indent    bool
}

//Format convert log to json text
func (f *JSONFormatter) Format(level, file, callstack string, msg interface{}) string {
	var jsonstr string
	f.Date = time.Now().Format("2006-01-02 15:04:05")
	f.Lev = level
	f.Msg = msg
	f.File = file
	if len(callstack) > 0 {
		f.Callstack = strings.Split(callstack, "\n")
		for i, line := range f.Callstack {
			f.Callstack[i] = strings.TrimLeft(line, "\t")
		}
	}

	data, err := json.Marshal(&f)
	if err != nil {
		jsonstr = fmt.Sprintf(`{"date":"%s","lev":"%s","error":"msg to json failed, %s"}`, f.Date, f.Lev, err.Error())
	} else {
		jsonstr = string(data)
	}

	if !f.indent {
		return jsonstr
	}

	var out bytes.Buffer
	err = json.Indent(&out, data, "", "\t")
	if err != nil {
		return jsonstr
	}

	return string(out.Bytes())
}
