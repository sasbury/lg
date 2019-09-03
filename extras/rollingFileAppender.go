package extras

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

/*
RollingFileAppender creates a series of log files.

A RollingFile appender will log to a file specified by prefix, which can contain a path, and a suffix, like "log". The appender will
concatenate the prefix and suffix using the following format "prefix.#.suffix" where # is the log file number. The current file will be "prefix.suffix".
Note, the . between the elements, the prefix and suffix should not include these.

Files can be rolled on size or manually by calling Roll().

The maxFiles must be at least 1
MaxFileSize must be at least 1024 - and is measured in bytes, if the max files is 1 the max file size is ignored

The actual file size will exceed maxFileSize, because the roller will not roll until a log message pushes the file past the size.
*/
type RollingFileAppender struct {
	sync.Mutex
	prefix        string
	suffix        string
	maxFileSize   int64
	maxFiles      int16
	firstTime     bool
	currentFile   *os.File
	currentWriter *bufio.Writer
}

// NewRollingFileAppender is used to create a rolling file appender.
func NewRollingFileAppender(prefix string, suffix string, maxFileSize int64, maxFiles int16) *RollingFileAppender {

	if maxFiles <= 0 {
		maxFiles = 1
	}

	if maxFileSize < 1024 {
		maxFileSize = 1024
	}

	appender := &RollingFileAppender{
		maxFileSize: maxFileSize,
		prefix:      prefix,
		suffix:      suffix,
		maxFiles:    maxFiles,
		firstTime:   true,
	}

	return appender
}

// currentFileName should be called inside the lock.
func (appender *RollingFileAppender) currentFileName() string {
	return fmt.Sprintf("%v.%v", appender.prefix, appender.suffix)
}

// Assumes the lock is held.
func (appender *RollingFileAppender) open() error {
	if appender.currentWriter != nil {
		return nil
	}

	f, err := os.OpenFile(appender.currentFileName(), os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		f, err = os.Create(appender.currentFileName())

		if err != nil {
			return err
		}
	}

	appender.currentFile = f
	appender.currentWriter = bufio.NewWriter(appender.currentFile)

	return nil
}

// Close closes the current file after flushing any buffered data.
// Locks the appender
func (appender *RollingFileAppender) Close() error {
	appender.Lock()
	defer appender.Unlock()
	return appender.close()
}

// close the current writer and file, assumes the lock is held
func (appender *RollingFileAppender) close() error {
	var err error

	if appender.currentWriter != nil {
		err = appender.currentWriter.Flush()
		appender.currentWriter = nil
	}

	if appender.currentFile != nil {
		err = appender.currentFile.Close()
		appender.currentFile = nil
	}

	return err
}

// needsRoll should be called inside the lock.
func (appender *RollingFileAppender) needsRoll() bool {
	if appender.maxFiles == 1 {
		_, err := os.Stat(appender.currentFileName())
		if err != nil {
			return os.IsNotExist(err)
		}
		return false
	}

	if appender.firstTime {
		return true
	}

	info, err := os.Stat(appender.currentFileName())

	if err != nil {
		return true
	}

	if info.Size() >= appender.maxFileSize {
		return true
	}

	return false
}

// Roll moves the file to the next number, up to the max files.
// locks the appender and closes the current file
func (appender *RollingFileAppender) Roll() error {
	appender.Lock()
	defer appender.Unlock()
	return appender.roll()
}

// assumes the lock is held
func (appender *RollingFileAppender) roll() error {
	appender.close()

	appender.firstTime = false

	for i := appender.maxFiles - 2; i >= 0; i-- {

		var fileName string

		if i == 0 {
			fileName = appender.currentFileName()
		} else {
			fileName = fmt.Sprintf("%v.%d.%v", appender.prefix, i, appender.suffix)
		}

		_, err := os.Stat(fileName)

		if err != nil {
			if os.IsNotExist(err) {
				continue // do'nt have this file yet
			} else {
				return err
			}
		}

		// we work backward so the only time the next file should exist is for the truly last file
		nextFileName := fmt.Sprintf("%v.%d.%v", appender.prefix, i+1, appender.suffix)
		_, err = os.Stat(nextFileName)

		if err != nil && !os.IsNotExist(err) {
			err = os.Remove(nextFileName)

			if err != nil {
				return err
			}
		}

		err = os.Rename(fileName, nextFileName)

		if err != nil {
			return err
		}
	}

	return nil
}

// Log a record to the current file.
func (appender *RollingFileAppender) Log(entry string) error {
	appender.Lock()
	defer appender.Unlock()

	if appender.needsRoll() {
		err := appender.roll()

		if err != nil {
			return err
		}

		err = appender.open()

		if err != nil {
			return err
		}
	}

	if appender.currentWriter == nil {
		err := appender.open()
		if err != nil {
			return err
		}
	}

	if appender.currentWriter != nil {
		_, err := appender.currentWriter.Write([]byte(entry))

		if err != nil {
			return err
		}

		_, err = appender.currentWriter.Write([]byte("\n"))

		if err != nil {
			return err
		}

		appender.currentWriter.Flush()
	}

	return nil
}
