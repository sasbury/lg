# lg

[![Build Status](https://travis-ci.org/sasbury/lg.svg?branch=master)](https://travis-ci.org/sasbury/lg) [![GoDoc](https://godoc.org/github.com/sasbury/lg?status.svg)](https://godoc.org/github.com/sasbury/lg)[![Go Report Card](https://goreportcard.com/badge/github.com/sasbury/lg)](https://goreportcard.com/report/github.com/sasbury/lg)

lg is a simple logging library loosely based on [Dave Cheney's Blog Post](https://dave.cheney.net/2015/11/05/lets-talk-about-logging) a few years ago. The fundamental idea is that most applications only need two log levels: info and debug. I have added the idea of tags which can be used to enable and disable debug mode on a finer grain.

Logging in lg is organized around the idea of a logger, string tags, a formatter function and an appender function.

```go
logger := lg.NewLogger()
```

Once you have a logger, you can simply print or debug with it:

```go
logger.Debugf("two %s", "formatted")
logger.Printf("four %s", "formatted")
```

Loggers are thread safe. But they do not protect their appender, see below, with their lock.

## Setting the Level

When calling debug, the formatting will happen after a debug flag is checked so there is no price for formatting or getting the time if the debug flag is false.

Debugging can be enabled for the entire logger:

```go
logger.EnableDebugMode()
logger.DisableDebugMode()

```

Debugging can also be enabled for specific tags:

```go
logger.EnableDebugModeFor("red")
logger.DisableDebugModeFor("red")
```

And you can turn off debugging for everything in one swoop:

```go
logger.DisableDebugModeAll()
```

Using tags for debugging does add a small price, if the global flag is true, this price is minimal, but if the global flag is false the library checks all the tags, until one is found with debugging enabled or the list is exhausted. This check is linear, scaling with the number of tags in the call and/or the debug tags.

## Formatting and Appenders

Loggers default to std err and the simple formatter. You can customize this with the Configure function:

```go
logger.Configure(lg.MinimalFormat, lg.StdErrAppender)
```

Custom formatters and appenders are supported:

```go
type LogFormatter func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string
type LogAppender func(entry string) error
```

The current release contains several formatters:

* `lg.FullFormat` - prints the time, log level and tags, if any
* `lg.SimpleFormat` - prints the time and log level
* `lg.MinimalFormat` - prints no extra data

This release also includes several appenders:

* `StdErrAppender` - writes to standard error, `os.Stderr`
* `StdOutAppender` - writes to standard out, `os.Stdout`
* `NullAppender` - no-op
* `ArrayAppender` - a struct that implements LogAppender, useful for tests.

## Loggers as Writers

Loggers implement the Writer interface, so you can attach them to the standard logging library:

```go
logger := NewLogger()
l, err := logger.Write([]byte("one formatted"))
stdLogger := log.New(logger, "", 0)
stdLogger.Print("three formatted")
```

In this situation you will want to either use the minimal formatting on the lg side, or 0 flags on the go logging side to avoid multiple headers/prefixes.

## Extras

The `extras` folder contains a few add-ons that aren't required but may be useful.

* RollingFileAppender - logs to a file and will roll based on size, with a max count of files
* BranchingAppender - appends to multiple child appenders
* BadAppender - always returns an error, useful for testing
* TestInjector - a logger based way to inject changes into production code for tests
