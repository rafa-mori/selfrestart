package logger

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	l "github.com/rafa-mori/logz"
)

type gLog struct {
	l.Logger
	gLogLevel LogType
}

var (
	// debug is a boolean that indicates whether to log debug messages.
	debug bool
	// g is the global logger instance.
	g *gLog = &gLog{
		Logger:    l.GetLogger("GoCrafter"),
		gLogLevel: LogTypeInfo,
	}
)

func init() {
	// Set the debug flag to true for testing purposes.
	debug = false
	// Initialize the global logger instance with a default logger.
	if g.Logger == nil {
		g = &gLog{
			Logger:    l.GetLogger("GoCrafter"),
			gLogLevel: LogTypeInfo,
		}
	}
}

type LogType string

const (
	LogTypeNotice  LogType = "notice"
	LogTypeInfo    LogType = "info"
	LogTypeDebug   LogType = "debug"
	LogTypeError   LogType = "error"
	LogTypeWarn    LogType = "warn"
	LogTypeFatal   LogType = "fatal"
	LogTypePanic   LogType = "panic"
	LogTypeSuccess LogType = "success"
)

// SetDebug is a function that sets the debug flag for logging.
func SetDebug(d bool) { debug = d }

// LogObjLogger is a function that logs messages with the specified log type.
func LogObjLogger[T any](obj *T, logType string, messages ...string) {
	if obj == nil {
		g.ErrorCtx(fmt.Sprintf("log object (%s) is nil", reflect.TypeFor[T]()), map[string]any{
			"context":  "Log",
			"logType":  logType,
			"object":   obj,
			"msg":      messages,
			"showData": true,
		})
		return
	}
	var lgr l.Logger
	if objValueLogger := reflect.ValueOf(obj).Elem().MethodByName("GetLogger"); !objValueLogger.IsValid() {
		if objValueLogger = reflect.ValueOf(obj).Elem().FieldByName("Logger"); !objValueLogger.IsValid() {
			g.ErrorCtx(fmt.Sprintf("log object (%s) does not have a logger field", reflect.TypeFor[T]()), map[string]any{
				"context":  "Log",
				"logType":  logType,
				"object":   obj,
				"msg":      messages,
				"showData": true,
			})
			return
		} else {
			lgrC := objValueLogger.Convert(reflect.TypeFor[l.Logger]())
			if lgrC.IsNil() {
				lgrC = reflect.ValueOf(g.Logger)
			}
			if lgr = lgrC.Interface().(l.Logger); lgr == nil {
				lgr = g.Logger
			}
		}
	} else {
		lgr = g.Logger
	}
	// First level of caller information is the Logger.go itself, so we skip it by using skip=2
	funcName, file, line, ok := getCallerInfo(2)
	if !ok {
		lgr.ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	ctxMessageMap := map[string]any{
		"context":  funcName,
		"file":     file,
		"line":     line,
		"showData": debug,
	}
	fullMessage := strings.Join(messages, " ")
	logType = strings.ToLower(logType)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(lgr, lType, fullMessage, ctxMessageMap)
		} else {
			lgr.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		lgr.InfoCtx(fullMessage, ctxMessageMap)
	}
}

// Log is a function that logs messages with the specified log type and caller information.
func Log(logType string, messages ...any) {
	// First level of caller information is the Logger.go itself, so we skip it by using skip=2
	funcName, file, line, ok := getCallerInfo(2)
	if !ok {
		g.ErrorCtx("Log: unable to get caller information", nil)
		return
	}
	ctxMessageMap := map[string]any{
		"context":  funcName,
		"file":     file,
		"line":     line,
		"showData": debug,
	}
	fullMessage := ""
	if len(messages) > 0 {
		fullMessage = fmt.Sprintf("%v", messages[0:])
	}
	logType = strings.ToLower(logType)
	if logType != "" {
		if reflect.TypeOf(logType).ConvertibleTo(reflect.TypeFor[LogType]()) {
			lType := LogType(logType)
			ctxMessageMap["logType"] = logType
			logging(g.Logger, lType, fullMessage, ctxMessageMap)
		} else {
			g.ErrorCtx(fmt.Sprintf("logType (%s) is not valid", logType), ctxMessageMap)
		}
	} else {
		g.InfoCtx(fullMessage, ctxMessageMap)
	}
}

// logging is a helper function that logs messages with the specified log type.
func logging(lgr l.Logger, lType LogType, fullMessage string, ctxMessageMap map[string]interface{}) {
	debugCtx := debug
	if !debugCtx {
		if lType == "error" || lType == "fatal" || lType == "panic" || lType == "debug" {
			// If debug is false, set the debug value based on the logType
			debugCtx = true
		} else {
			debugCtx = false
		}
	}
	ctxMessageMap["showData"] = debugCtx
	switch lType {
	case LogTypeInfo:
		fmt.Print("  ")
		lgr.InfoCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeDebug:
		lgr.DebugCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeError:
		lgr.ErrorCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeWarn:
		lgr.WarnCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeNotice:
		lgr.NoticeCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeSuccess:
		lgr.SuccessCtx(fullMessage, ctxMessageMap)
		break
	case LogTypeFatal:
		lgr.FatalCtx(fullMessage, ctxMessageMap)
		break
	case LogTypePanic:
		lgr.FatalCtx(fullMessage, ctxMessageMap)
		break
	default:
		lgr.InfoCtx(fullMessage, ctxMessageMap)
		break
	}
	debugCtx = debug
}

// getCallerInfo extracts caller information and formats it appropriately
func getCallerInfo(skip int) (string, string, int, bool) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", 0, false
	}

	funcName := runtime.FuncForPC(pc).Name()
	if funcName == "" {
		funcName = "unknown"
	}

	// Simplify function name (remove package path)
	if lastSlash := strings.LastIndex(funcName, "/"); lastSlash >= 0 {
		if lastSlash+1 < len(funcName) {
			funcName = funcName[lastSlash+1:]
		} else {
			funcName = "unknown"
		}
	}

	// Format file path based on debug mode
	formattedFile := ""

	if debug {
		formattedFile += file
	} else {
		fileParts := strings.Split(filepath.Clean(filepath.Base(file)), "/")
		if len(fileParts) > 0 {
			// Use the last part of the file path
			formattedFile = fileParts[len(fileParts)-1]
		} else {
			formattedFile = filepath.Clean(filepath.Base(file))
		}
		// Extract just the filename without extension
		fileName := filepath.Base(file)
		if ext := filepath.Ext(fileName); ext != "" {
			// Remove function name extension if it exists
			fileName = strings.TrimSuffix(fileName, ext)
		}
	}
	if debug {
		return formattedFile, formattedFile, line, ok
	}
	timeStampTime := time.Now().Format("2006-01-02 15:04:05")
	timeStampStr := fmt.Sprintf("%s", timeStampTime)
	return timeStampStr, "", line, ok
}
