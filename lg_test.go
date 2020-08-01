package lg

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNilAppenderForCoverage(t *testing.T) {
	count := 0
	logger := NewLogger()
	logger.Configure(func(debug bool, tags []string, t time.Time, fmt string, args ...interface{}) string {
		count++
		return ""
	}, nil)

	logger.EnableDebugMode()
	logger.Printf("test %s", "test")
	logger.Debugf("test %s", "test")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 0, count)
}

func TestNilLogger(t *testing.T) {
	var logger *Logger

	err := logger.Printf("test %s", "test")
	require.NoError(t, err)

	err = logger.Debugf("test %s", "test")
	require.NoError(t, err)

	err = logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.NoError(t, err)

	err = logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")
	require.NoError(t, err)
}

func TestFullFormatDebugOff(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(FullFormat, a.Log)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 2, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "three formatted"))
	require.True(t, strings.Contains(a.Entries[1], "four formatted"))
	require.True(t, strings.Contains(a.Entries[0], "[INF]"))
	require.True(t, strings.Contains(a.Entries[1], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.True(t, strings.Contains(a.Entries[1], "red"))
	require.True(t, strings.Contains(a.Entries[1], "blue"))
}

func TestFullFormatDebugOn(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(FullFormat, a.Log)
	logger.EnableDebugMode()

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 4, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "one formatted"))
	require.True(t, strings.Contains(a.Entries[1], "two formatted"))
	require.True(t, strings.Contains(a.Entries[2], "three formatted"))
	require.True(t, strings.Contains(a.Entries[3], "four formatted"))
	require.True(t, strings.Contains(a.Entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.Entries[1], "[DBG]"))
	require.True(t, strings.Contains(a.Entries[2], "[INF]"))
	require.True(t, strings.Contains(a.Entries[3], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.True(t, strings.Contains(a.Entries[1], "red"))
	require.True(t, strings.Contains(a.Entries[1], "blue"))
	require.False(t, strings.Contains(a.Entries[2], "red"))
	require.False(t, strings.Contains(a.Entries[2], "blue"))
	require.True(t, strings.Contains(a.Entries[3], "red"))
	require.True(t, strings.Contains(a.Entries[3], "blue"))
}

func TestSimpleFormatDebugOff(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(SimpleFormat, a.Log)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 2, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "three formatted"))
	require.True(t, strings.Contains(a.Entries[1], "four formatted"))
	require.True(t, strings.Contains(a.Entries[0], "[INF]"))
	require.True(t, strings.Contains(a.Entries[1], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.False(t, strings.Contains(a.Entries[1], "red"))
	require.False(t, strings.Contains(a.Entries[1], "blue"))
}

func TestSimpleFormatDebugOn(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(SimpleFormat, a.Log)
	logger.EnableDebugMode()

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 4, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "one formatted"))
	require.True(t, strings.Contains(a.Entries[1], "two formatted"))
	require.True(t, strings.Contains(a.Entries[2], "three formatted"))
	require.True(t, strings.Contains(a.Entries[3], "four formatted"))
	require.True(t, strings.Contains(a.Entries[0], "[DBG]"))
	require.True(t, strings.Contains(a.Entries[1], "[DBG]"))
	require.True(t, strings.Contains(a.Entries[2], "[INF]"))
	require.True(t, strings.Contains(a.Entries[3], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.False(t, strings.Contains(a.Entries[1], "red"))
	require.False(t, strings.Contains(a.Entries[1], "blue"))
	require.False(t, strings.Contains(a.Entries[2], "red"))
	require.False(t, strings.Contains(a.Entries[2], "blue"))
	require.False(t, strings.Contains(a.Entries[3], "red"))
	require.False(t, strings.Contains(a.Entries[3], "blue"))
}

func TestMinimalFormatDebugOff(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(MinimalFormat, a.Log)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 2, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "three formatted"))
	require.True(t, strings.Contains(a.Entries[1], "four formatted"))
	require.False(t, strings.Contains(a.Entries[0], "[INF]"))
	require.False(t, strings.Contains(a.Entries[1], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "[INF]"))
	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.False(t, strings.Contains(a.Entries[1], "red"))
	require.False(t, strings.Contains(a.Entries[1], "blue"))
}

func TestMinimalFormatDebugOn(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(MinimalFormat, a.Log)
	logger.EnableDebugMode()

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")

	require.Equal(t, 4, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "one formatted"))
	require.True(t, strings.Contains(a.Entries[1], "two formatted"))
	require.True(t, strings.Contains(a.Entries[2], "three formatted"))
	require.True(t, strings.Contains(a.Entries[3], "four formatted"))
	require.False(t, strings.Contains(a.Entries[0], "[DBG]"))
	require.False(t, strings.Contains(a.Entries[1], "[DBG]"))
	require.False(t, strings.Contains(a.Entries[2], "[INF]"))
	require.False(t, strings.Contains(a.Entries[3], "[INF]"))

	require.False(t, strings.Contains(a.Entries[0], "red"))
	require.False(t, strings.Contains(a.Entries[0], "blue"))
	require.False(t, strings.Contains(a.Entries[1], "red"))
	require.False(t, strings.Contains(a.Entries[1], "blue"))
	require.False(t, strings.Contains(a.Entries[2], "red"))
	require.False(t, strings.Contains(a.Entries[2], "blue"))
	require.False(t, strings.Contains(a.Entries[3], "red"))
	require.False(t, strings.Contains(a.Entries[3], "blue"))
}

func TestDebugStatus(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(MinimalFormat, a.Log)

	require.False(t, logger.IsDebugMode())
	require.False(t, logger.IsDebugModeFor("red"))
	require.False(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))

	logger.EnableDebugModeFor("red")
	require.False(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))
	require.False(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))

	logger.EnableDebugModeFor("blue")
	require.False(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))
	require.True(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))

	logger.EnableDebugMode()
	require.True(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))
	require.True(t, logger.IsDebugModeFor("blue"))
	require.True(t, logger.IsDebugModeFor("green"))

	logger.DisableDebugMode()
	require.False(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))
	require.True(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))

	logger.DisableDebugModeFor("red")
	require.False(t, logger.IsDebugMode())
	require.False(t, logger.IsDebugModeFor("red"))
	require.True(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))

	logger.DisableDebugModeAll()
	require.False(t, logger.IsDebugMode())
	require.False(t, logger.IsDebugModeFor("red"))
	require.False(t, logger.IsDebugModeFor("blue"))
	require.False(t, logger.IsDebugModeFor("green"))
}

func TestDebugFlags(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(MinimalFormat, a.Log)

	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 0, len(a.Entries))

	logger.EnableDebugMode()
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 1, len(a.Entries))
	require.True(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))

	logger.DisableDebugMode()
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 1, len(a.Entries))
	require.False(t, logger.IsDebugMode())
	require.False(t, logger.IsDebugModeFor("red"))

	logger.EnableDebugModeFor("red")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 2, len(a.Entries))
	require.False(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("red"))

	logger.DisableDebugModeFor("red")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 2, len(a.Entries))
	require.False(t, logger.IsDebugMode())
	require.False(t, logger.IsDebugModeFor("red"))

	logger.EnableDebugModeFor("green")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	require.Equal(t, 2, len(a.Entries))
	require.False(t, logger.IsDebugMode())
	require.True(t, logger.IsDebugModeFor("green"))

	logger.EnableDebugModeFor("blue")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.TagDebugf([]string{"yellow"}, "two %s", "formatted") // shouldn't be logged
	require.Equal(t, 3, len(a.Entries))

	logger.EnableDebugMode()
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.TagDebugf([]string{"yellow"}, "two %s", "formatted") // should log with global flag
	require.Equal(t, 5, len(a.Entries))

	logger.DisableDebugMode()
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.TagDebugf([]string{"yellow"}, "two %s", "formatted") // should no longer log with global flag
	require.Equal(t, 6, len(a.Entries))

	logger.DisableDebugModeAll()
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.TagDebugf([]string{"yellow"}, "two %s", "formatted")
	require.Equal(t, 6, len(a.Entries))
}

func TestStdOutAppender(t *testing.T) { // For coverage
	logger := NewLogger()
	logger.Configure(MinimalFormat, StdOutAppender)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")
}

func TestStdErrAppender(t *testing.T) { // For coverage
	logger := NewLogger()
	logger.Configure(MinimalFormat, StdErrAppender)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")
}

func TestNullAppender(t *testing.T) { // For coverage
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)

	logger.Debugf("one %s", "formatted")
	logger.TagDebugf([]string{"red", "blue"}, "two %s", "formatted")
	logger.Printf("three %s", "formatted")
	logger.TagPrintf([]string{"red", "blue"}, "four %s", "formatted")
}

func TestWriterImplementation(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(FullFormat, a.Log)

	l, err := logger.Write([]byte("one formatted"))
	require.NoError(t, err)
	require.Len(t, "one formatted", l)
	l, err = logger.Write([]byte("two formatted"))
	require.NoError(t, err)
	require.Len(t, "two formatted", l)

	// Test adapted log framework
	stdLogger := log.New(logger, "", 0)
	stdLogger.Print("three formatted")
	stdLogger.Printf("four %s", "formatted")

	require.Equal(t, 4, len(a.Entries))
	require.True(t, strings.Contains(a.Entries[0], "one formatted"))
	require.True(t, strings.Contains(a.Entries[1], "two formatted"))
	require.True(t, strings.Contains(a.Entries[2], "three formatted"))
	require.True(t, strings.Contains(a.Entries[3], "four formatted"))
	require.True(t, strings.Contains(a.Entries[0], "[INF]"))
	require.True(t, strings.Contains(a.Entries[1], "[INF]"))
	require.True(t, strings.Contains(a.Entries[2], "[INF]"))
	require.True(t, strings.Contains(a.Entries[3], "[INF]"))
}

func BenchmarkDebugWithDebugOff(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)

	for n := 0; n < b.N; n++ {
		logger.Debugf("one %s", "formatted")
	}
}

func BenchmarkDebugWithDebugOffWithTags(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)

	tags := []string{"red", "blue"}
	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}

func BenchmarkDebugWithDebugOn(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugMode()

	for n := 0; n < b.N; n++ {
		logger.Debugf("one %s", "formatted")
	}
}

func BenchmarkTagDebugFirst(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("red", "green", "blue")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}

func BenchmarkTagDebugSecond(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("blue", "red", "green")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}

func BenchmarkTagDebugThird(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("green", "blue", "red")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}

func BenchmarkTagDebugWithDebugOn(b *testing.B) {
	b.ReportAllocs()
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("red")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}
