lg
================

[![Build Status](https://travis-ci.org/sasbury/lg.svg?branch=master)](https://travis-ci.org/sasbury/lg) [![GoDoc](https://godoc.org/github.com/sasbury/lg?status.svg)](https://godoc.org/github.com/sasbury/lg)

lg is a simple logging library loosely based on [Dave Cheney's Blog Post](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) a few years ago. The fundamental idea is that most applications only need two log levels: info and debug.

Logging in lg is organized around the idea of a logger. Loggers are named, and a default logger is available:

```go
logger := lg.GetLogger("my-logger")
logger2 := lg.DefaultLogger()
```

Using an empty string, "", for the name will result in the default logger.

Once you have a logger, you can simply print or debug with it:

```go
logger.Debug("one")
logger.Debugf("two %s", "formatted")
logger.Print("three")
logger.Printf("four %s", "formatted")
```

When calling debug, the formatting will happen after a debug flag is checked.

Loggers can be configured using a single function:

```go
lg.ConfigLogger(loggerName, false, []string{"red", "blue"}, lg.SimpleFormat, lg.StdErrAppender)
```

config logger takes 5 arguments:

* The name of the logger to configure
* Whether or not to enable debug logging
* An array of tags, these may or may not be displayed depending on the format.
* A format function of the form `type LogFormatter func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string`
* A logging appender function of the form `type LogAppender func(entry string)`

Tags can be used to help searching based on the logger output.

The current release contains several formatters:

* `lg.FullFormat` - prints the time, log level and tags, if any
* `lg.SimpleFormat` - prints the time and log level
* `lg.MinimalFormat` - prints no extra data

This release also includes several appenders:

* `StdErrAppender` - writes to standard error, `os.Stderr`
* `StdOutAppender` - writes to standard out, `os.Stdout`
* `NullAppender` - no-op
* `ArrayAppender` - a struct that implements LogAppender, useful for tests.
