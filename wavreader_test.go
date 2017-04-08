package audiotransport

import (
	"os"
	"testing"

	"crypto/rand"

	"github.com/gurupras/gocommons"
	"github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {
	require := require.New(t)

	file, err := gocommons.Open("test.wav", os.O_RDONLY, gocommons.GZ_FALSE)
	require.Nil(err, "Failed to open file", err)

	// First read in the file headers
	fh, err := ParseFileHeader(file.File)
	require.Nil(err, "Failed to parse file header", err)
	require.Equal("RIFF", string(fh.ID[:]))
	require.Equal(902212, fh.Length)
	require.Equal("WAVE", string(fh.Type[:]))

	ch, err := ParseChunkHeader(file.File)
	require.Nil(err, "Failed to parse chunk header", err)
	require.Equal("fmt ", string(ch.ID[:]))
	require.Equal(16, ch.Length)

	fmtHeader, err := ParseFormatHeader(file.File)
	require.Nil(err, "Failed to parse format header", err)

	require.Equal(uint16(1), fmtHeader.AudioFormat)
	require.Equal(uint16(2), fmtHeader.NumChannels)
	require.Equal(44100, fmtHeader.SampleRate)
	require.Equal(176400, fmtHeader.ByteRate)
	require.Equal(uint16(4), fmtHeader.BlockAlign)
	require.Equal(uint16(16), fmtHeader.BitsPerSample)

	ch, err = ParseChunkHeader(file.File)
	require.Nil(err, "Failed to parse chunk header", err)
	require.Equal("data", string(ch.ID[:]))
	require.Equal(902176, ch.Length)
}

func TestPlayback(t *testing.T) {
	require := require.New(t)

	file, err := gocommons.Open("test.wav", os.O_RDONLY, gocommons.GZ_FALSE)
	require.Nil(err, "Failed to open file", err)

	wh, err := ParseWavHeaders(file.File)
	require.Nil(err, "Failed to parse headers")
	_ = wh

	dataChunkHeader, err := ParseChunkHeader(file.File)
	require.Nil(err, "Failed to parse chunk header", err)
	require.Equal("data", string(dataChunkHeader.ID[:]))

	data := make([]byte, dataChunkHeader.Length)
	n, err := file.File.Read(data)
	require.Nil(err, "Failed to read data from WAV file", err)
	require.Equal(dataChunkHeader.Length, n)

	backend := &AlsaBackend{}
	backend.Init("TestPlayback", "default", uint32(wh.SampleRate), uint32(wh.NumChannels), true)
	_, err = backend.Write(data)
	require.Nil(err)
}

func RandomBytes(b []byte, len int) {
	_, _ = rand.Read(b)
}
func TestRandomPlayback(t *testing.T) {
	require := require.New(t)

	backend := &AlsaBackend{}
	backend.Init("TestRandomPlayback", "default", 44100, 2, true)

	bufsize := 512 * 1024
	buf := make([]byte, bufsize)
	RandomBytes(buf, bufsize)

	_, err := backend.Write(buf)
	require.Nil(err)
}
