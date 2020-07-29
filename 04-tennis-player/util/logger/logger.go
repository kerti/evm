package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"

	syslog "github.com/RackSec/srslog"
	"golang.org/x/net/context"
)

const (
	// LogLevelTrace outputs TRACE, DEBUG, INFO, WARNING, ERROR, and PANIC logs
	LogLevelTrace = "TRACE"
	// LogLevelDebug outputs DEBUG, INFO, WARNING, ERROR, and PANIC logs
	LogLevelDebug = "DEBUG"
	// LogLevelInfo outputs INFO, WARNING, ERROR, and PANIC logs
	LogLevelInfo = "INFO"
	// LogLevelWarning outputs WARNING, ERROR, and PANIC logs
	LogLevelWarning = "WARNING"
	// LogLevelError outputs ERROR and PANIC logs
	LogLevelError = "ERROR"
	// LogLevelPanic outputs PANIC logs only
	LogLevelPanic = "PANIC"
)

var logLevel = map[string]int{
	LogLevelTrace:   5,
	LogLevelDebug:   4,
	LogLevelInfo:    3,
	LogLevelWarning: 2,
	LogLevelError:   1,
	LogLevelPanic:   0,
}

var logColor = map[string]string{
	LogLevelTrace:   "90",
	LogLevelDebug:   "34",
	LogLevelInfo:    "32",
	LogLevelWarning: "33",
	LogLevelError:   "31",
	LogLevelPanic:   "91",
}

var logPrefix = map[string]string{
	LogLevelTrace:   fmt.Sprintf("\033[1m\033[%smTRACE: \033[0m", logColor[LogLevelTrace]),
	LogLevelDebug:   fmt.Sprintf("\033[1m\033[%smDEBUG: \033[0m", logColor[LogLevelDebug]),
	LogLevelInfo:    fmt.Sprintf("\033[1m\033[%smINFO : \033[0m", logColor[LogLevelInfo]),
	LogLevelWarning: fmt.Sprintf("\033[1m\033[%smWARN : \033[0m", logColor[LogLevelWarning]),
	LogLevelError:   fmt.Sprintf("\033[1m\033[%smERROR: \033[0m", logColor[LogLevelError]),
	LogLevelPanic:   fmt.Sprintf("\033[1m\033[%smPANIC: \033[0m", logColor[LogLevelPanic]),
}

var ptSystemName string

var activeLogLevel = strings.ToUpper(os.Getenv("LOG_LEVEL"))

func parseLogLevel() string {
	switch activeLogLevel {
	case LogLevelTrace, LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelPanic:
	default:
		activeLogLevel = LogLevelTrace
	}
	return activeLogLevel
}

func getActiveLogLevel() int {
	return logLevel[activeLogLevel]
}

func extractReqID(ctx context.Context) string {
	requestIDKey := "x-request-id"
	str, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return str
}

// SetupLogger creates logger instance to log to PaperTrail and Console. Should only called once in main function.
func SetupLogger(ptHost string, ptPort string) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if ptHost != "" {
		hostname, _ := os.Hostname()
		ptEndpoint := fmt.Sprintf("%s:%s", ptHost, ptPort)
		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatal("Can't connect to PaperTrail ...")
		}

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("No papertrail transport detected. Logger only use local stdout")
	}
}

func formatterRFC3164(p syslog.Priority, hostname, tag, content string) string {
	return syslog.RFC3164Formatter(p, ptSystemName, hostname, content)
}

// SetupLoggerAuto creates logger instance to log to PaperTrail automatically without specifying PT HOST and PORT
func SetupLoggerAuto(appName string, ptEndpoint string) {
	activeLogLevel = parseLogLevel()
	log.SetPrefix("")
	log.SetFlags(0)

	if appName != "" && ptEndpoint != "" {
		ptSystemName = appName
		hostname, _ := os.Hostname()

		ptWriter, err := syslog.Dial("udp", ptEndpoint, syslog.LOG_INFO, hostname)

		if err != nil {
			log.Fatalf("Can't connect to PaperTrail: %s", err.Error())
		}

		ptWriter.SetFormatter(formatterRFC3164)

		log.SetOutput(io.MultiWriter(os.Stdout, ptWriter))
	} else {
		log.Print("Logger configured to use only local stdout")
	}
}

// Warn prints warning message to logs
func Warn(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelWarning] {
		message := fmt.Sprintf(logPrefix[LogLevelWarning]+format, v...)
		log.Print(message)
	}
}

// Trace prints trace message to logs
func Trace(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelTrace] {
		message := fmt.Sprintf(logPrefix[LogLevelTrace]+format, v...)
		log.Print(message)
	}
}

// Debug prints debug message to logs
func Debug(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelDebug] {
		message := fmt.Sprintf(logPrefix[LogLevelDebug]+format, v...)
		log.Print(message)
	}
}

// Info prints info message to logs
func Info(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelInfo] {
		message := fmt.Sprintf(logPrefix[LogLevelInfo]+format, v...)
		log.Print(message)
	}
}

// Err prints error message to logs
func Err(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelError]+format, v...)}
		message = append(message, "\n", string(debug.Stack()))
		log.Print(message...)
	}
}

// ErrNoStack prints error message to logs without stacktrace
func ErrNoStack(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelError]+format, v...)}
		log.Print(message...)
	}
}

// Panic prints panic message to logs
func Panic(format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelPanic] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelPanic]+format, v...)}
		log.Print(message...)
	}
}

// Fatal calls Err and then os.Exit(1)
func Fatal(format string, v ...interface{}) {
	Err(format, v...)
	os.Exit(1)
}

// WarnContext prints warning message to logs
func WarnContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelWarning] {
		message := fmt.Sprintf(logPrefix[LogLevelWarning]+" ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// TraceContext prints trace message to logs
func TraceContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelTrace] {
		message := fmt.Sprintf(logPrefix[LogLevelTrace]+" ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// DebugContext prints debug message to logs
func DebugContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelDebug] {
		message := fmt.Sprintf(logPrefix[LogLevelDebug]+" ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// InfoContext prints info message to logs
func InfoContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelInfo] {
		message := fmt.Sprintf(logPrefix[LogLevelInfo]+" ReqID "+extractReqID(ctx)+" - "+format, v...)
		log.Print(message)
	}
}

// ErrContext prints error message to logs
func ErrContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelError]+" ReqID "+extractReqID(ctx)+" - "+format, v...)}
		message = append(message, "\n", string(debug.Stack()))
		log.Print(message...)
	}
}

// ErrNoStackContext prints error message to logs without stacktrace
func ErrNoStackContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelError] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelError]+" ReqID "+extractReqID(ctx)+" - "+format, v...)}
		log.Print(message...)
	}
}

// PanicContext prints panic message to logs
func PanicContext(ctx context.Context, format string, v ...interface{}) {
	if getActiveLogLevel() >= logLevel[LogLevelPanic] {
		message := []interface{}{fmt.Sprintf(logPrefix[LogLevelPanic]+" ReqID "+extractReqID(ctx)+" - "+format, v...)}
		log.Print(message...)
	}
}

// FatalContext calls Err and then os.Exit(1)
func FatalContext(ctx context.Context, format string, v ...interface{}) {
	ErrContext(ctx, format, v...)
	os.Exit(1)
}
