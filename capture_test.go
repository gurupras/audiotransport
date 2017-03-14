package audiotransport

import (
	"os"
	"sync"
	"testing"

	"github.com/gurupras/audiotransport/alsa"
	"github.com/gurupras/gocommons"
	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	assert := assert.New(t)

	isPlaying := true
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		idx := alsa.Pa_init("TestCapture", "alsa_output.pci-0000_00_05.0.analog-stereo.monitor", 48000, 2, 0)
		assert.True(idx >= 0, "Failed to initialize pulseaudio", idx)
		defer alsa.Pa_release(idx)

		size := 48000
		buf := make([]byte, size)

		file, err := gocommons.Open("capture.wav", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, gocommons.GZ_FALSE)
		assert.Nil(err, "Failed to open capture file", err)

		for {
			alsa.Pa_handle_read(idx, &buf, int32(size))
			file.File.Write(buf)
			mutex.Lock()
			if !isPlaying {
				break
			}
			mutex.Unlock()
		}
	}()

	go func() {
		playIdx := alsa.Pa_init("TestPlayback", "NULL", 48000, 2, 1)
		assert.True(playIdx >= 0, "Failed to initialize pulseaudio", playIdx)

		size := 48000
		buf := make([]byte, size)

		file, err := gocommons.Open("test.wav", os.O_RDONLY, gocommons.GZ_FALSE)
		assert.Nil(err, "Failed to open WAV file", err)
		defer file.Close()

		_, err = ParseWavHeaders(file.File)
		assert.Nil(err, "Failed to parse WAV headers", err)

		dataHeader, err := ParseChunkHeader(file.File)
		assert.Nil(err, "Failed to parse data header", err)

		count := 0
		for count < dataHeader.Length {
			remaining := dataHeader.Length - count
			bufsize := size
			if remaining < size {
				bufsize = remaining
			}
			file.File.Read(buf)
			ret := alsa.Pa_handle_write(playIdx, &buf, int32(bufsize))
			assert.Equal(int32(0), ret, "pa_handle_write failed")
			count += size
		}
		mutex.Lock()
		isPlaying = false
		mutex.Unlock()

		//FIXME: This line crashes the test
		//alsa.Pa_release(playIdx)
	}()

	wg.Wait()
}
