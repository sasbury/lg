package extras

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sasbury/lg"
)

// TestInjector is a log appender that tracks a map of patterns to TestInjections. When a patterns the test injector is discovered in
// a log entry the TestInjection is executed and the result is returned.  This pattern can be used to inject errors into
// test code that is normally hard to force, like causing a file not found after creating a file.
//
// TestInjector also demonstrates how a LogAppender can work like middleware and call a "next" appender.
//
// As an example, suppose you have server code like:
//
// func (ms *server) doWork() (string, error) {
// 	ms.logger.TagDebugf([]string{"server", "stat"}, "checking file %s", ms.filepath)
// 	info, err := os.Stat(ms.filepath)
// 	if os.IsNotExist(err) {
// 		return "", err
// 	}
//
// 	ms.logger.TagPrintf([]string{"server", "dircheck"}, "checking if %s is a directory", ms.filepath)
// 	if info.IsDir() {
// 		return "", fmt.Errorf("%s is a directory", ms.filepath)
// 	}
//
// 	ms.logger.TagDebugf([]string{"server", "readfile"}, "reading data at %s", ms.filepath)
// 	content, err := ioutil.ReadFile(ms.filepath)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	return string(content), nil
// }
//
// You can inject an error for the file read using something like:
//
// injector.Add("readfile", func() error {
//		return os.Rename(filepath, newFilePath)
//	})
//
// where the file path is the one used by the server.
//
// Because the injector is implemented as a log appender, tags are found with string matching, which means that very specific
// tags should be used. Also, the injector is not "fast", and scales linearly with the number of injections. Then again, most
// use cases can be handled with 1 injection, by setting the logger being called once for each test.
//
// The injector requires that the format used includes tags! If you want to hide tags you could create a similar pattern with
// a formatter which will see the tags earlier. Also, if the message won't print, ie is debug and debug is off, then the injector
// won't be executed.
type TestInjector struct {
	sync.Mutex
	next      lg.LogAppender
	callbacks map[string]TestInjection
}

// TestInjection is a simple function that the TestInjector tracks and executes when it sees a specific tag.
type TestInjection func() error

// InjectorError holds a list of child errors
type InjectorError struct {
	Children []error
}

func (be InjectorError) Error() string {
	return fmt.Sprintf("injector error with %d children", len(be.Children))
}

// NewTestInjector returns a new test injector with the provided next appender.
func NewTestInjector(next lg.LogAppender) *TestInjector {
	return &TestInjector{
		next:      next,
		callbacks: map[string]TestInjection{},
	}
}

// Add a tag + TestInjection pair to the injector
func (inj *TestInjector) Add(pattern string, callback TestInjection) {
	inj.Lock()
	inj.callbacks[pattern] = callback
	inj.Unlock()
}

// Remove a tag + TestInjection pair to the injector
func (inj *TestInjector) Remove(pattern string) {
	inj.Lock()
	delete(inj.callbacks, pattern)
	inj.Unlock()
}

// Clear all injectors
func (inj *TestInjector) Clear() {
	inj.Lock()
	inj.callbacks = map[string]TestInjection{}
	inj.Unlock()
}

// Log the entry, if any errors are encountered either in the next appender or in
// a callback, return an injector error. Callbacks are not run inside the lock, they are
// copied outside it to prevent possible deadlocks, however this can create weird situations
// if injectors are added/removed during logging
func (inj *TestInjector) Log(entry string) error {
	var errors []error

	toRun := []TestInjection{}

	// lock and read but call everything outside the lock
	inj.Lock()
	for t, cb := range inj.callbacks {
		if strings.Contains(entry, t) && cb != nil {
			toRun = append(toRun, cb)
		}
	}
	nextLogger := inj.next
	inj.Unlock()

	for _, cb := range toRun {
		err := cb()
		if err != nil {
			errors = append(errors, err)
		}
	}

	if nextLogger != nil {
		err := nextLogger(entry)

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return BranchingError{
			Children: errors,
		}
	}

	return nil
}
