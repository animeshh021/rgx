package log

import (
	"fmt"
	"os"
	"time"
)

type loglevels int

const (
	trace loglevels = iota
	debug
	info
)

func isTerminal() bool {
	fileInfo, e := os.Stdout.Stat()
	if e != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func canShowColor() bool {
	// see https://no-color.org/
	return isTerminal() && os.Getenv("RGX_COLOR_LOGS") == "1"
}

//goland:noinspection GoUnusedConst
const (
	fgBlack   = 30
	fgRed     = 31
	fgGreen   = 32
	fgYellow  = 33
	fgBlue    = 34
	fgMagenta = 35
	fgCyan    = 36
	fgWhite   = 37
	fgDefault = 38
	fgReset   = 0
)

var logLevel = info
var colorizeOutput = canShowColor()

const tsFormat = "15:04:05"

func EnableDebug() {
	logLevel = debug
}

func EnableTrace() {
	logLevel = trace
}

const ansiReset = "\x1b[0;0m"

func color(fg int, bright bool) string {
	fgmode := 0
	if bright {
		fgmode = 1
	}
	return fmt.Sprintf("\x1b[%d;%dm", fgmode, fg)
}

func toColor(level string) string {
	switch level {
	case "TRACE":
		return color(fgWhite, false)
	case "DEBUG":
		return color(fgCyan, false)
	case "INFO":
		return color(fgGreen, false)
	case "WARN":
		return color(fgYellow, false)
	case "ERROR":
		return color(fgRed, false)
	case "FATAL":
		return color(fgMagenta, true)
	}
	return color(fgDefault, false)
}

func log(level, msg string) {
	if colorizeOutput {
		fmt.Printf("%s[%s %s] %s%s\n", toColor(level), time.Now().Format(tsFormat), level, msg, ansiReset)
	} else {
		fmt.Printf("[%s %s] %s\n", time.Now().Format(tsFormat), level, msg)
	}

}

func Trace(v ...any) {
	//goland:noinspection GoBoolExpressions
	if logLevel == trace {
		log("TRACE", fmt.Sprintf(v[0].(string), v[1:]...))
	}
}

func Debug(v ...any) {
	if logLevel <= debug {
		log("DEBUG", fmt.Sprintf(v[0].(string), v[1:]...))
	}
}

func Info(v ...any) {
	log("INFO", fmt.Sprintf(v[0].(string), v[1:]...))
}

func Warn(v ...any) {
	log("WARN", fmt.Sprintf(v[0].(string), v[1:]...))
}

func Error(v ...any) {
	log("ERROR", fmt.Sprintf(v[0].(string), v[1:]...))
}

func Fatal(v ...any) {
	log("FATAL", fmt.Sprintf(v[0].(string), v[1:]...))
	os.Exit(11)
}
