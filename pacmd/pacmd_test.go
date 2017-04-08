package pacmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListSources(t *testing.T) {
	require := require.New(t)

	sources, err := ListSources()
	require.Nil(err)
	require.True(len(sources) > 0)
}

func TestListSinks(t *testing.T) {
	require := require.New(t)

	sinks, err := ListSinks()
	require.Nil(err)
	require.True(len(sinks) > 0)
}
