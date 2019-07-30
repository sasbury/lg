package lg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNamedLoggers(t *testing.T) {
	logger := GetLogger("named-logger")
	logger2 := GetLogger("named-logger")

	require.True(t, logger == logger2, "named loggers should be the same")
}

func TestDefaultName(t *testing.T) {
	logger := GetLogger("")
	logger2 := DefaultLogger()
	require.True(t, logger == logger2, "empty string should go to default logger")
}

func TestFullFormatDebugOff(t *testing.T) {
	loggerName := "test-full-format-off"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, false, []string{"red", "blue"}, FullFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 2, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "three"))
	require.True(t, strings.Contains(a.entries[0], "[INF]"))
	require.True(t, strings.Contains(a.entries[1], "four formatted"))

	require.True(t, strings.Contains(a.entries[0], "red"))
	require.True(t, strings.Contains(a.entries[0], "blue"))
}

func TestFullFormatDebugOn(t *testing.T) {
	loggerName := "test-full-format-on"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, true, []string{"red", "blue"}, FullFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "one"))
	require.True(t, strings.Contains(a.entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.entries[1], "two formatted"))
	require.True(t, strings.Contains(a.entries[2], "three"))
	require.True(t, strings.Contains(a.entries[2], "[INF]"))
	require.True(t, strings.Contains(a.entries[3], "four formatted"))

	require.True(t, strings.Contains(a.entries[0], "red"))
	require.True(t, strings.Contains(a.entries[0], "blue"))
}

func TestFullFormatNoTags(t *testing.T) {
	loggerName := "test-full-format-notags"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, true, nil, FullFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "one"))
	require.True(t, strings.Contains(a.entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.entries[3], "four formatted"))
	require.True(t, strings.Contains(a.entries[3], "[INF]"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}

func TestSimpleFormatDebugOff(t *testing.T) {
	loggerName := "test-simple-format-off"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, false, []string{"red", "blue"}, SimpleFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 2, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "three"))
	require.True(t, strings.Contains(a.entries[0], "[INF]"))
	require.True(t, strings.Contains(a.entries[1], "four formatted"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}

func TestSimpleFormatDebugOn(t *testing.T) {
	loggerName := "test-simple-format-on"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, true, []string{"red", "blue"}, SimpleFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "one"))
	require.True(t, strings.Contains(a.entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.entries[1], "two formatted"))
	require.True(t, strings.Contains(a.entries[2], "three"))
	require.True(t, strings.Contains(a.entries[2], "[INF]"))
	require.True(t, strings.Contains(a.entries[3], "four formatted"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}

func TestMinimalFormatDebugOff(t *testing.T) {
	loggerName := "test-minimal-format-off"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, false, []string{"red", "blue"}, MinimalFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 2, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "three"))
	require.False(t, strings.Contains(a.entries[0], "[INF]")) // no debug flag in minimal
	require.True(t, strings.Contains(a.entries[1], "four formatted"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}

func TestMinimalFormatDebugOn(t *testing.T) {
	loggerName := "test-minimal-format-on"
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, true, []string{"red", "blue"}, MinimalFormat, a.log)
	require.NoError(t, err)

	logger := GetLogger(loggerName)

	logger.Debug("one")
	logger.Debugf("two %s", "formatted")
	logger.Print("three")
	logger.Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "one"))
	require.False(t, strings.Contains(a.entries[0], "[DBG]")) // no debug flag in minimal
	require.True(t, strings.Contains(a.entries[1], "two formatted"))
	require.True(t, strings.Contains(a.entries[2], "three"))
	require.False(t, strings.Contains(a.entries[2], "[INF]")) // no debug flag in minimal
	require.True(t, strings.Contains(a.entries[3], "four formatted"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}

func TestDefaultLogger(t *testing.T) {
	loggerName := ""
	a := &ArrayAppender{}

	err := ConfigLogger(loggerName, true, []string{"red", "blue"}, SimpleFormat, a.log)
	require.NoError(t, err)

	Debug("one")
	Debugf("two %s", "formatted")
	Print("three")
	Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.entries))
	require.True(t, strings.Contains(a.entries[0], "one"))
	require.True(t, strings.Contains(a.entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.entries[1], "two formatted"))
	require.True(t, strings.Contains(a.entries[2], "three"))
	require.True(t, strings.Contains(a.entries[2], "[INF]"))
	require.True(t, strings.Contains(a.entries[3], "four formatted"))

	require.False(t, strings.Contains(a.entries[0], "red"))
	require.False(t, strings.Contains(a.entries[0], "blue"))
}
