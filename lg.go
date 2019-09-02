package lg

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Logger provides a minimal configuration and methods to print and debug, with or without tags.
// Debugging can be configured at the tag level or for the entire logger.
type Logger struct {
	sync.RWMutex
	debug     bool
	debugTags map[string]bool
	format    LogFormatter
	appender  LogAppender
}

// LogFormatter is used to convert a logmessage to a string for printing
// The logger's lock will be used to protect the formatter
type LogFormatter func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string

// LogAppender is used to write the log message to a destination
// A new line is not assumed in the entry and should be added by the appender if appropriate
// The logger's lock will not be used protect the appender
type LogAppender func(entry string)

// NewLogger creates and returns a new logger
func NewLogger() *Logger {
	return &Logger{
		format:    SimpleFormat,
		appender:  StdErrAppender,
		debug:     false,
		debugTags: map[string]bool{},
	}
}

// EnableDebugMode turns on debug mode for all tags
func (l *Logger) EnableDebugMode() {
	l.Lock()
	l.debug = true
	l.Unlock()
}

// DisableDebugMode turns off the debug flag, individual tags may still have debug mode on
func (l *Logger) DisableDebugMode() {
	l.Lock()
	l.debug = false
	l.Unlock()
}

// DisableDebugModeAll turns off the debug flag, and removes any debug flags
func (l *Logger) DisableDebugModeAll() {
	l.Lock()
	l.debug = false
	l.debugTags = map[string]bool{}
	l.Unlock()
}

// IsDebugMode returns true if the debug flag is on
func (l *Logger) IsDebugMode() bool {
	l.Lock()
	debugMode := l.debug
	l.Unlock()
	return debugMode
}

// EnableDebugModeFor turn on debug mode for one or more tags
func (l *Logger) EnableDebugModeFor(tags ...string) {
	l.Lock()
	for _, t := range tags {
		l.debugTags[t] = true
	}
	l.Unlock()
}

// DisableDebugModeFor turns off debug mode for one or more tags
func (l *Logger) DisableDebugModeFor(tags ...string) {
	l.Lock()
	for _, t := range tags {
		delete(l.debugTags, t)
	}
	l.Unlock()
}

// IsDebugModeFor returns true if the debug flag is on for a specific tag
func (l *Logger) IsDebugModeFor(tag string) bool {
	l.Lock()
	debugMode := l.debug
	if !debugMode {
		tf, ok := l.debugTags[tag]
		debugMode = ok && tf
	}
	l.Unlock()
	return debugMode
}

// Configure set the formatter and appender
func (l *Logger) Configure(formatter LogFormatter, appender LogAppender) {
	l.Lock()
	l.format = formatter
	l.appender = appender
	l.Unlock()
}

//Printf used for most logging, prints the formatted string with the configured formatter
func (l *Logger) Printf(fmt string, args ...interface{}) {
	l.RLock()
	entry := l.format(false, nil, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	app(entry)
}

//Debugf prints the formatted string with the configured formatter, if debug is on
func (l *Logger) Debugf(fmt string, args ...interface{}) {
	l.RLock()
	if !l.debug {
		l.RUnlock()
		return
	}
	entry := l.format(true, nil, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	app(entry)
}

//TagPrintf used for most logging, prints the formatted string with the configured formatter
func (l *Logger) TagPrintf(tags []string, fmt string, args ...interface{}) {
	l.RLock()
	entry := l.format(false, tags, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	app(entry)
}

//TagDebugf prints the formatted string with the configured formatter, if debug mode is on for any of the tags
func (l *Logger) TagDebugf(tags []string, fmt string, args ...interface{}) {
	l.RLock()
	debugMode := l.debug
	if !debugMode {
		for _, t := range tags {
			tf, ok := l.debugTags[t]
			if ok && tf {
				debugMode = true
				break
			}
		}
	}
	if !debugMode {
		l.RUnlock()
		return
	}
	entry := l.format(true, tags, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	app(entry)
}

// FullFormat includes everything
func FullFormat(debug bool, tags []string, t time.Time, format string, args ...interface{}) string {
	formatStr := ""
	timeStr := t.Format(time.StampMilli)
	modeStr := "[INF]"

	if debug {
		modeStr = "[DBG]"
	}

	if tags != nil && len(tags) > 0 {
		formatStr = fmt.Sprintf("%s %s [%s] %s", timeStr, modeStr, strings.Join(tags, ", "), format)
	} else {
		formatStr = fmt.Sprintf("%s %s %s", timeStr, modeStr, format)
	}

	return fmt.Sprintf(formatStr, args...)
}

// SimpleFormat includes time, debug and message, but not the tags
func SimpleFormat(debug bool, tags []string, t time.Time, format string, args ...interface{}) string {
	formatStr := ""
	timeStr := t.Format(time.StampMilli)
	modeStr := "[INF]"

	if debug {
		modeStr = "[DBG]"
	}

	formatStr = fmt.Sprintf("%s %s %s", timeStr, modeStr, format)

	return fmt.Sprintf(formatStr, args...)
}

// MinimalFormat just formats the message, tags and time are ignored
func MinimalFormat(debug bool, tags []string, t time.Time, format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// StdErrAppender is an appender for stderr
func StdErrAppender(entry string) {
	fmt.Fprintln(os.Stderr, entry)
}

// StdOutAppender is an appender for stdout
func StdOutAppender(entry string) {
	fmt.Fprintln(os.Stdout, entry)
}

// NullAppender simply ignores all append calls
func NullAppender(entry string) {
}

// ArrayAppender stores entries in an array for testing
type ArrayAppender struct {
	Entries []string
}

func (a *ArrayAppender) log(entry string) {
	a.Entries = append(a.Entries, entry)
}
