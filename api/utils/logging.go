package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var (
	// DefaultLogger ...
	DefaultLogger *logrus.Logger

	// TrackingLogger ...
	TrackingLogger *logrus.Logger
	goPath         = os.Getenv("GOPATH")
	spewConfig     = spew.ConfigState{MaxDepth: 5, Indent: " ", DisableMethods: true}
)

const (
	stackTraceSize = 10
)

func init() {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05-0700"
	formatter.FullTimestamp = true

	DefaultLogger = logrus.New()
	DefaultLogger.Formatter = formatter
	DefaultLogger.Out = os.Stdout
	DefaultLogger.Level = logrus.DebugLevel

	jsonFormatter := new(logrus.JSONFormatter)
	jsonFormatter.TimestampFormat = "2006-01-02T15:04:05-0700"
	TrackingLogger = logrus.New()
	TrackingLogger.Formatter = jsonFormatter
	TrackingLogger.Out = os.Stdout
	TrackingLogger.Level = logrus.InfoLevel
}

// SetLogLevel ...
func SetLogLevel(logLevel string) {
	if DefaultLogger == nil {
		return
	}

	switch logLevel {
	case "DEBUG":
		DefaultLogger.Level = logrus.DebugLevel
	case "INFO":
		DefaultLogger.Level = logrus.InfoLevel
	case "WARN":
		DefaultLogger.Level = logrus.WarnLevel
	case "ERROR":
		DefaultLogger.Level = logrus.ErrorLevel
	case "FATAL":
		DefaultLogger.Level = logrus.FatalLevel
	}
}

// Info ...
func Info(tag string, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(nil)).Info(spewConfig.Sprintf(formatter, args...))
}

// InfoWithParams ...
func InfoWithParams(tag string, params map[string]interface{}, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(params)).Info(spewConfig.Sprintf(formatter, args...))
}

// Warn ...
func Warn(tag string, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(nil)).Warn(spewConfig.Sprintf(formatter, args...))
}

// WarnWithParams ...
func WarnWithParams(tag string, params map[string]interface{}, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(params)).Warn(spewConfig.Sprintf(formatter, args...))
}

// Debug ...
func Debug(tag string, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(nil)).Debug(spewConfig.Sprintf(formatter, args...))
}

// DebugWithParams ...
func DebugWithParams(tag string, params map[string]interface{}, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(params)).Debug(spewConfig.Sprintf(formatter, args...))
}

// Error ...
func Error(tag string, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	for _, arg := range args {
		if err, ok := arg.(error); ok && errors.Is(err, context.Canceled) {
			Warn(tag, format, args)
			return
		}
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(nil)).Error(spewConfig.Sprintf(formatter, args...))
}

// ErrorWithParams ...
func ErrorWithParams(tag string, params map[string]interface{}, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	for _, arg := range args {
		if err, ok := arg.(error); ok && errors.Is(err, context.Canceled) {
			WarnWithParams(tag, params, format, args)
			return
		}
	}

	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(params)).Error(spewConfig.Sprintf(formatter, args...))
}

// Track ...
func Track(tag string, args interface{}) {
	if TrackingLogger == nil {
		return
	}

	params := map[string]interface{}{"tag": tag}
	TrackingLogger.WithFields(addFileParams(params)).Info(getJSONFormat(args))
}

// TrackWithParams ...
func TrackWithParams(tag string, params map[string]interface{}, args interface{}) {
	if TrackingLogger == nil {
		return
	}

	params["tag"] = tag
	TrackingLogger.WithFields(addFileParams(params)).Info(getJSONFormat(args))
}

// Fatal ...
func Fatal(tag string, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	argList := getArgListWithTrace(args)
	format += "%v"
	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(nil)).Fatal(spewConfig.Sprintf(formatter, argList...))
}

// FatalWithParams ...
func FatalWithParams(tag string, params map[string]interface{}, format string, args ...interface{}) {
	if DefaultLogger == nil {
		return
	}

	argList := getArgListWithTrace(args)
	format += "%v"
	formatter := GetFormatter(tag, format)
	DefaultLogger.WithFields(addFileParams(params)).Fatal(spewConfig.Sprintf(formatter, argList...))
}

func GetFormatter(tag string, format string) string {
	formatter := fmt.Sprintf("%s : ", tag)
	formatter = strings.TrimPrefix(formatter, " : ")
	formatter += format
	return formatter
}

func addFileParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		params = map[string]interface{}{}
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	params["filename"] = fmt.Sprintf("%s:%d", file, line)

	return params
}

func getJSONFormat(args interface{}) string {
	var b []byte
	var err error
	if m, ok := args.(json.Marshaler); ok {
		b, err = m.MarshalJSON()
	} else {
		b, err = json.Marshal(args)
	}
	if err != nil {
		Error("Logging", "unable to marshall get json format, err=%v", err)
	}
	return string(b)
}

func getArgListWithTrace(args ...interface{}) []interface{} {
	argList := make([]interface{}, len(args))
	for k, v := range args {
		argList[k] = v
	}
	argList = append(argList, getTrace())

	return argList
}

func getTrace() string {
	pc := make([]uintptr, stackTraceSize)
	returnedSize := runtime.Callers(3, pc)

	buf := bytes.NewBuffer(nil)
	_, _ = buf.WriteString("\n-------- [STACK_TRACE] --------\n")
	for i := 0; i < returnedSize; i++ {
		pc[i]--
		f := runtime.FuncForPC(pc[i] + 1)
		file, line := f.FileLine(pc[i])
		if relPath, err := filepath.Rel(goPath, file); err == nil {
			if !strings.HasPrefix(relPath, "..") {
				file = relPath
			}
		}
		_, _ = buf.WriteString(fmt.Sprintf("%s:%d\n\t%s\n", file, line, f.Name()))
	}
	_, _ = buf.WriteString("--------------------------------\n")

	return string(buf.Bytes())
}
