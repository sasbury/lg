package lg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullFormatDebugOff(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(FullFormat, a.log)

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
	logger.Configure(FullFormat, a.log)
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
	logger.Configure(SimpleFormat, a.log)

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
	logger.Configure(SimpleFormat, a.log)
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
	logger.Configure(MinimalFormat, a.log)

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
	logger.Configure(MinimalFormat, a.log)
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

func TestDebugFlags(t *testing.T) {
	a := &ArrayAppender{}

	logger := NewLogger()
	logger.Configure(MinimalFormat, a.log)

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

func BenchmarkDebugWithDebugOff(b *testing.B) {
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)

	for n := 0; n < b.N; n++ {
		logger.Debugf("one %s", "formatted")
	}
}

func BenchmarkDebugWithDebugOn(b *testing.B) {
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugMode()

	for n := 0; n < b.N; n++ {
		logger.Debugf("one %s", "formatted")
	}
}

func BenchmarkTagDebugWithDebugOff(b *testing.B) {
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("red")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}

func BenchmarkTagDebugWithDebugOn(b *testing.B) {
	logger := NewLogger()
	logger.Configure(MinimalFormat, NullAppender)
	logger.EnableDebugModeFor("red")
	tags := []string{"red"}

	for n := 0; n < b.N; n++ {
		logger.TagDebugf(tags, "one %s", "formatted")
	}
}
