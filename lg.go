package lg

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Log is the interface for the objects that are the target of logging messages. Logging methods
// imply a level. For example, Info() implies a level of LogLevel.INFO.
type Log interface {
	Printf(fmt string, args ...interface{})
	Print(args ...interface{})

	Debugf(fmt string, args ...interface{})
	Debug(args ...interface{})
}

// LogFormatter is used to convert a logmessage to a string for printing
type LogFormatter func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string

// LogAppender is used to write the log message to a destination
// A new line is not assumed in the entry and should be added by the appender if appropriate
type LogAppender func(entry string)

var loggers = make(map[string]*logger)
var loggersMutex = sync.RWMutex{}

type logger struct {
	sync.RWMutex
	name     string
	debug    bool
	format   LogFormatter
	tags     []string
	appender LogAppender
}

// DefaultLogger returns a logger that can be used when a named logger isn't required.
func DefaultLogger() Log {
	return GetLogger("_default")
}

// GetLogger returns a named logger, creating it if necessary.
// Use ConfigLogger to assign settings
func GetLogger(name string) Log {
	loggersMutex.Lock()

	if name == "" {
		name = "_default"
	}

	log := loggers[name]

	if log == nil {
		log = &logger{}
		log.name = name
		log.debug = false
		log.format = SimpleFormat
		log.appender = StdErrAppender
		loggers[name] = log
	}
	loggersMutex.Unlock()

	return log
}

// ConfigLogger creates and/or configures the underlying logger
// An empty name will target the default logger
func ConfigLogger(name string, debug bool, tags []string, formatter LogFormatter, appender LogAppender) error {
	loggersMutex.Lock()

	if name == "" {
		name = "_default"
	}

	log := loggers[name]

	if log == nil {
		log = &logger{}
		log.name = name
		log.debug = false
		loggers[name] = log
	}
	loggersMutex.Unlock()

	log.Lock()
	defer log.Unlock()
	log.debug = debug
	log.tags = tags
	log.format = formatter
	log.appender = appender

	return nil
}

//Printf used for most logging, prints the formatted string with the configured formatter
func (l *logger) Printf(fmt string, args ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	entry := l.format(false, l.tags, time.Now(), fmt, args...)
	l.appender(entry)
}

//Print prints the formatted string with the configured formatter
func (l *logger) Print(args ...interface{}) {
	l.RLock()
	defer l.RUnlock()
	entry := l.format(false, l.tags, time.Now(), fmt.Sprint(args...))
	l.appender(entry)
}

//Debugf prints the formatted string with the configured formatter, if debug is on
func (l *logger) Debugf(fmt string, args ...interface{}) {
	l.RLock()
	defer l.RUnlock()

	if !l.debug {
		return
	}
	entry := l.format(true, l.tags, time.Now(), fmt, args...)
	l.appender(entry)
}

//Debug prints the formatted string with the configured formatter, if debug is on
func (l *logger) Debug(args ...interface{}) {
	l.RLock()
	defer l.RUnlock()

	if !l.debug {
		return
	}
	entry := l.format(true, l.tags, time.Now(), fmt.Sprint(args...))
	l.appender(entry)
}

// Printf using the default logger
func Printf(fmt string, args ...interface{}) {
	DefaultLogger().Printf(fmt, args...)
}

// Print using the default logger
func Print(args ...interface{}) {
	DefaultLogger().Print(args...)
}

// Debugf using the default logger
func Debugf(fmt string, args ...interface{}) {
	DefaultLogger().Debugf(fmt, args...)
}

// Debug using the default logger
func Debug(args ...interface{}) {
	DefaultLogger().Debug(args...)
}

// FullFormat includes everything
func FullFormat(debug bool, tags []string, t time.Time, format string, args ...interface{}) string {
	formatStr := ""
	timeStr := t.Format(time.StampMilli)

	if debug {
		if tags != nil && len(tags) > 0 {
			formatStr = fmt.Sprintf("%s [DBG] [%s] %s", timeStr, strings.Join(tags, ", "), format)
		} else {
			formatStr = fmt.Sprintf("%s [DBG] %s", timeStr, format)
		}
	} else {
		if tags != nil && len(tags) > 0 {
			formatStr = fmt.Sprintf("%s [INF] [%s] %s", timeStr, strings.Join(tags, ", "), format)
		} else {
			formatStr = fmt.Sprintf("%s [INF] %s", timeStr, format)
		}
	}

	return fmt.Sprintf(formatStr, args...)
}

// SimpleFormat includes time, debug and message, but not the tags
func SimpleFormat(debug bool, tags []string, t time.Time, format string, args ...interface{}) string {
	formatStr := ""
	timeStr := t.Format(time.StampMilli)

	if debug {
		formatStr = fmt.Sprintf("%s [DBG] %s", timeStr, format)
	} else {
		formatStr = fmt.Sprintf("%s [INF] %s", timeStr, format)
	}

	return fmt.Sprintf(formatStr, args...)
}

// MinimalFormat just formats the message
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
