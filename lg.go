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
	debugTags []string
	format    LogFormatter
	appender  LogAppender
}

// LogFormatter is used to convert a logmessage to a string for printing
// The logger's lock will be used to protect the formatter
type LogFormatter func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string

// LogAppender is used to write the log message to a destination
// A new line is not assumed in the entry and should be added by the appender if appropriate
// The logger's lock will not be used protect the appender
type LogAppender func(entry string) error

// NewLogger creates and returns a new default logger
func NewLogger() *Logger {
	return &Logger{
		format:    SimpleFormat,
		appender:  StdErrAppender,
		debug:     false,
		debugTags: []string{},
	}
}

// NewLoggerWithConfig creates and returns a new logger with the specified format and appender
func NewLoggerWithConfig(formatter LogFormatter, appender LogAppender) *Logger {
	return &Logger{
		format:    formatter,
		appender:  appender,
		debug:     false,
		debugTags: []string{},
	}
}

// Write implements the Writer interface, so that a logger can be used to adapt
// the standard go.log library. In this case you probably want to set the flags
// to 0 for go logging, or use the minimal appender on the lg side.
// This is equivalent to calling "Printf"
func (l *Logger) Write(p []byte) (n int, err error) {
	s := string(p[:])
	l.Printf(s)
	return len(p), nil
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
	l.debugTags = []string{}
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
		l.debugTags = append(l.debugTags, t)
	}
	l.Unlock()
}

// DisableDebugModeFor turns off debug mode for one or more tags
func (l *Logger) DisableDebugModeFor(tags ...string) {
	l.Lock()
	tagSet := map[string]bool{}
	for _, t := range l.debugTags {
		tagSet[t] = true
	}
	for _, t := range tags {
		delete(tagSet, t)
	}
	newTags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		newTags = append(newTags, t)
	}
	l.debugTags = newTags
	l.Unlock()
}

// IsDebugModeFor returns true if the debug flag is on for a specific tag
func (l *Logger) IsDebugModeFor(tag string) bool {
	l.Lock()
	debugMode := l.debug
	if !debugMode {
		ld := len(l.debugTags)
		for i := 0; i < ld; i++ {
			if l.debugTags[i] == tag {
				debugMode = true
				break
			}
		}
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
//A nil logger will do nothing
func (l *Logger) Printf(fmt string, args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.RLock()
	if l.appender == nil {
		l.RUnlock()
		return nil
	}
	entry := l.format(false, nil, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	return app(entry)
}

//Debugf prints the formatted string with the configured formatter, if debug is on
func (l *Logger) Debugf(fmt string, args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.RLock()
	if !l.debug || l.appender == nil {
		l.RUnlock()
		return nil
	}
	entry := l.format(true, nil, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	return app(entry)
}

//TagPrintf used for most logging, prints the formatted string with the configured formatter
func (l *Logger) TagPrintf(tags []string, fmt string, args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.RLock()
	if l.appender == nil {
		l.RUnlock()
		return nil
	}
	entry := l.format(false, tags, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	return app(entry)
}

//TagDebugf prints the formatted string with the configured formatter, if debug mode is on for any of the tags
func (l *Logger) TagDebugf(tags []string, fmt string, args ...interface{}) error {
	if l == nil {
		return nil
	}
	l.RLock()
	if l.appender == nil {
		l.RUnlock()
		return nil
	}
	debugMode := l.debug
	if !debugMode {
		ld, lt := len(l.debugTags), len(tags)

		for i := 0; i < ld; i++ {
			for j := 0; j < lt; j++ {
				if l.debugTags[i] == tags[j] {
					debugMode = true
					break
				}
			}
			if debugMode {
				break
			}
		}
	}
	if !debugMode {
		l.RUnlock()
		return nil
	}
	entry := l.format(true, tags, time.Now(), fmt, args...)
	app := l.appender
	l.RUnlock()
	return app(entry)
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
func StdErrAppender(entry string) error {
	fmt.Fprintln(os.Stderr, entry)
	return nil
}

// StdOutAppender is an appender for stdout
func StdOutAppender(entry string) error {
	fmt.Fprintln(os.Stdout, entry)
	return nil
}

// NullAppender simply ignores all append calls
func NullAppender(entry string) error {
	return nil
}

// ArrayAppender stores entries in an array for testing
type ArrayAppender struct {
	sync.Mutex
	Entries []string
}

// Log is ArrayAppenders implementation of LogAppender
func (a *ArrayAppender) Log(entry string) error {
	a.Lock()
	a.Entries = append(a.Entries, entry)
	a.Unlock()
	return nil
}
