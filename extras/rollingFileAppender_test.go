package extras

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/sasbury/lg"
	"github.com/stretchr/testify/require"
)

func TestRollingAppender(t *testing.T) {
	filepath := path.Join(os.TempDir(), "appendtest")
	lfAppender := NewRollingFileAppender(filepath, "log", int64(2048), 5)
	arrayAppender := &lg.ArrayAppender{}
	appender := NewBranchingAppender(lfAppender.Log, arrayAppender.Log)

	logger := lg.NewLogger()
	logger.Configure(lg.MinimalFormat, appender.Log)

	for i := 0; i < 2548; i++ {
		logger.Printf("1")
		logger.Debugf("1")
	}

	lfAppender.Close()

	require.Equal(t, len(arrayAppender.Entries), 2548, "should have logged all the messages")

	pathOne := fmt.Sprintf("%s.log", filepath)
	info, err := os.Stat(pathOne)
	require.Nil(t, err, "Stat should be able to find the log file")
	require.Equal(t, info.Size(), int64(1000), "new file should have 1000 bytes, 500 1's and new lines")

	pathTwo := fmt.Sprintf("%s.1.log", filepath)
	info, err = os.Stat(pathTwo)
	require.Nil(t, err, "Stat should be able to find the rolled log file")
	require.Equal(t, info.Size(), int64(2048), "rolled file should have 2048 bytes")
}

func TestRollingAppenderOneFile(t *testing.T) {
	entries := 2548
	filepath := path.Join(os.TempDir(), "appendtest")
	lfAppender := NewRollingFileAppender(filepath, "log", int64(2048), 1)
	arrayAppender := &lg.ArrayAppender{}
	appender := NewBranchingAppender(lfAppender.Log, arrayAppender.Log)

	logger := lg.NewLogger()
	logger.Configure(lg.MinimalFormat, appender.Log)

	pathOne := fmt.Sprintf("%s.log", filepath)
	err := os.Remove(pathOne) //we won't roll so make sure we start with nothing

	if err != nil && !os.IsNotExist(err) {
		require.Nil(t, err, "Should be able to delete")
	}

	for i := 0; i < entries; i++ {
		logger.Printf("1")
		logger.Debugf("1") // none of these will log
	}

	lfAppender.Close()

	require.Equal(t, len(arrayAppender.Entries), entries, "should have logged all the messages")

	pathOne = fmt.Sprintf("%s.log", filepath)
	info, err := os.Stat(pathOne)
	require.Nil(t, err, "Stat should be able to find the log file")
	require.Equal(t, info.Size(), int64(entries*2), "new file should have all the data, since there isn't any rolling")
}

func TestRollingAppenderNew(t *testing.T) {
	filepath := path.Join(os.TempDir(), "appendtest")
	app := NewRollingFileAppender(filepath, "log", int64(100), -1)

	require.Equal(t, app.maxFiles, int16(1), "max files defaults to 1")
	require.Equal(t, app.maxFileSize, int64(1024), "max filesize defaults to 1024")
	require.Equal(t, app.currentFileName(), fmt.Sprintf("%s.%s", filepath, "log"), "current file name is always prefix.suffix")
}
