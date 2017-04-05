package audiotransport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPulseSources(t *testing.T) {
	require := require.New(t)

	sources, err := PulseSources()
	require.Nil(err)
	require.True(len(sources) > 0)
}

func TestPulseSinks(t *testing.T) {
	require := require.New(t)

	sinks, err := PulseSinks()
	require.Nil(err)
	require.True(len(sinks) > 0)
}
