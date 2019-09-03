package extras

import (
	"strings"
	"testing"

	"github.com/sasbury/lg"
	"github.com/stretchr/testify/require"
)

func TestBranchingAppender(t *testing.T) {
	app1 := &lg.ArrayAppender{}
	app2 := &lg.ArrayAppender{}
	appender := NewBranchingAppender(app1.Log, app2.Log)

	appender.Log("one")
	appender.Log("two")
	appender.Log("three")

	require.Len(t, app1.Entries, 3)
	require.Len(t, app2.Entries, 3)

	for i, v := range app1.Entries {
		require.Equal(t, v, app2.Entries[i])
	}
}

func TestBranchingErrors(t *testing.T) {
	app1 := &lg.ArrayAppender{}
	appender := NewBranchingAppender(app1.Log, BadAppender)

	err := appender.Log("one")
	require.Error(t, err)
	require.Len(t, err.(BranchingError).Children, 1)

	err = appender.Log("two")
	require.Error(t, err)
	require.Len(t, err.(BranchingError).Children, 1)

	// for coverage
	require.True(t, strings.Contains(err.Error(), "1"))

	err = appender.Log("three")
	require.Error(t, err)
	require.Len(t, err.(BranchingError).Children, 1)

	require.Len(t, app1.Entries, 3)
}
