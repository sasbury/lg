package extras

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sasbury/lg"
	"github.com/stretchr/testify/require"
)

type server struct {
	logger   *lg.Logger
	filepath string
}

func (ms *server) doWork() (string, error) {
	ms.logger.TagDebugf([]string{"server", "stat"}, "checking file %s", ms.filepath)
	info, err := os.Stat(ms.filepath)
	if os.IsNotExist(err) {
		return "", err
	}

	ms.logger.TagPrintf([]string{"server", "dircheck"}, "checking if %s is a directory", ms.filepath)
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", ms.filepath)
	}

	ms.logger.TagDebugf([]string{"server", "readfile"}, "reading data at %s", ms.filepath)
	content, err := ioutil.ReadFile(ms.filepath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func TestInjectorAtDifferentPoints(t *testing.T) {
	expected := "hello world"
	f, err := ioutil.TempFile(os.TempDir(), "injector")
	require.NoError(t, err)

	dirpath, err := ioutil.TempDir(os.TempDir(), "inject")
	require.NoError(t, err)

	filepath := f.Name()
	ioutil.WriteFile(filepath, []byte(expected), 0666)

	arrApp := &lg.ArrayAppender{}
	injector := NewTestInjector(arrApp.Log)
	logger := lg.NewLogger()
	logger.Configure(lg.FullFormat, injector.Log)
	logger.EnableDebugModeFor("server")

	ms := &server{
		logger:   logger,
		filepath: filepath,
	}

	// Run it once normally
	content, err := ms.doWork()
	require.NoError(t, err)
	require.Equal(t, expected, content)

	// Make the stat fail
	newFilePath := filepath + ".1"
	injector.Add("stat", func() error {
		return os.Rename(filepath, newFilePath)
	})

	content, err = ms.doWork() // should fail with not exist
	require.Error(t, err)
	require.Empty(t, content)
	require.True(t, os.IsNotExist(err))

	// Fix stat, but make the dir check fail, this could be done without the injector
	injector.Clear()
	ioutil.WriteFile(filepath, []byte(expected), 0666)
	injector.Add("dircheck", func() error {
		ms.filepath = dirpath
		return nil
	})

	content, err = ms.doWork() // should fail with not exist
	require.Error(t, err)
	require.Empty(t, content)
	require.False(t, os.IsNotExist(err))

	// Fix stat and make the read fail
	ms.filepath = filepath
	injector.Add("readfile", func() error {
		return os.Rename(filepath, newFilePath)
	})

	content, err = ms.doWork() // should fail with not exist
	require.Error(t, err)
	require.Empty(t, content)
	require.False(t, os.IsNotExist(err))
}
