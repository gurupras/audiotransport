package audiotransport

import (
	"os"
	"sync"
	"testing"

	"crypto/rand"

	"github.com/gurupras/audiotransport/alsa"
	"github.com/gurupras/gocommons"
	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	assert := assert.New(t)

	file, err := gocommons.Open("test.wav", os.O_RDONLY, gocommons.GZ_FALSE)
	assert.Nil(err, "Failed to open file", err)

	// First read in the file headers
	fh, err := ParseFileHeader(file.File)
	assert.Nil(err, "Failed to parse file header", err)
	assert.Equal("RIFF", string(fh.ID[:]))
	assert.Equal(902212, fh.Length)
	assert.Equal("WAVE", string(fh.Type[:]))

	ch, err := ParseChunkHeader(file.File)
	assert.Nil(err, "Failed to parse chunk header", err)
	assert.Equal("fmt ", string(ch.ID[:]))
	assert.Equal(16, ch.Length)

	fmtHeader, err := ParseFormatHeader(file.File)
	assert.Nil(err, "Failed to parse format header", err)

	assert.Equal(uint16(1), fmtHeader.AudioFormat)
	assert.Equal(uint16(2), fmtHeader.NumChannels)
	assert.Equal(44100, fmtHeader.SampleRate)
	assert.Equal(176400, fmtHeader.ByteRate)
	assert.Equal(uint16(4), fmtHeader.BlockAlign)
	assert.Equal(uint16(16), fmtHeader.BitsPerSample)

	ch, err = ParseChunkHeader(file.File)
	assert.Nil(err, "Failed to parse chunk header", err)
	assert.Equal("data", string(ch.ID[:]))
	assert.Equal(902176, ch.Length)
}

func TestPlayback(t *testing.T) {
	assert := assert.New(t)

	file, err := gocommons.Open("test.wav", os.O_RDONLY, gocommons.GZ_FALSE)
	assert.Nil(err, "Failed to open file", err)

	wh, err := ParseWavHeaders(file.File)
	assert.Nil(err, "Failed to parse headers")
	_ = wh

	dataChunkHeader, err := ParseChunkHeader(file.File)
	assert.Nil(err, "Failed to parse chunk header", err)
	assert.Equal("data", string(dataChunkHeader.ID[:]))

	data := make([]byte, dataChunkHeader.Length)
	n, err := file.File.Read(data)
	assert.Nil(err, "Failed to read data from WAV file", err)
	assert.Equal(dataChunkHeader.Length, n)

	idx := alsa.Alsa_init("default", int32(wh.SampleRate), int32(wh.NumChannels), 1)
	alsa.Alsa_play_bytes(idx, &data, int32(len(data)))
}

func RandomBytes(b []byte, len int) {
	_, _ = rand.Read(b)
}
func TestRandomPlayback(t *testing.T) {
	idx := alsa.Alsa_init("default", 44100, 2, 1)
	// One thread to generate and one thread to play
	bufsize := 512 * 1024
	var buf *[]byte
	buf1 := make([]byte, bufsize)
	RandomBytes(buf1, bufsize)
	buf = &buf1

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for {
			alsa.Alsa_play_bytes(idx, buf, int32(bufsize))
		}
	}()
	wg.Wait()
}
