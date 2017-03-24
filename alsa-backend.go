package audiotransport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/gurupras/goalsa"
)

type AlsaBackend struct {
	*Backend
	AlsaInterface
}

type AlsaInterface interface {
	io.Reader
	io.Writer
}
type AlsaPlaybackDevice struct {
	*alsa.PlaybackDevice
}

func (playback *AlsaPlaybackDevice) Read(buf []byte) (n int, err error) {
	return 0, errors.New("ALSA playback device cannot perform Read()")
}

func (playback *AlsaPlaybackDevice) Write(buf []byte) (samples int, err error) {
	var data []interface{}
	data = convert(buf, playback.FormatSampleSize(), playback.Endian())
	return playback.PlaybackDevice.Write(data)
}

type AlsaCaptureDevice struct {
	*alsa.CaptureDevice
}

func (capture *AlsaCaptureDevice) Read(buf []byte) (samples int, err error) {
	var data []interface{}
	data = convert(buf, capture.FormatSampleSize(), capture.Endian())
	// Now ship data off to goalsa
	return capture.CaptureDevice.Read(data)
}

func (capture *AlsaCaptureDevice) Write(buf []byte) (n int, err error) {
	return 0, errors.New("ALSA capture device cannot perform Write()")
}

func convert(buf []byte, size int, endian binary.ByteOrder) (data []interface{}) {
	chunk := make([]byte, size)
	var entry interface{}

	if size == 1 {
		data = make([]interface{}, len(buf))
		for i := 0; i < len(buf); i++ {
			data[i] = buf[i]
		}
		return
	} else {
		data = make([]interface{}, 0)
	}

	//FIXME: This should not be hard-coded to binary.LittleEndian
	// We should be using the endianness that the device has been configured with
	buffer := bytes.NewBuffer(buf)
	for i := 0; i < len(buf); i += size {
		buffer.Read(chunk)
		switch size {
		case 2:
			entry = endian.Uint16(chunk)
		case 4:
			entry = endian.Uint32(chunk)
		case 8:
			entry = endian.Uint64(chunk)
		}
		data = append(data, entry)
	}
	return
}

func (ab *AlsaBackend) Init(name, device string, samplerate, channels uint32, isPlayback bool) (err error) {
	if ab.Backend == nil {
		ab.Backend = &Backend{}
	}
	ab.Backend.Init(name, device, samplerate, channels)

	var alsaInterface AlsaInterface

	if isPlayback {
		var dev *alsa.PlaybackDevice
		if dev, err = alsa.NewPlaybackDevice(device, int(channels), alsa.FormatS16LE, int(samplerate), alsa.BufferParams{}); err != nil {
			return
		}
		alsaInterface = &AlsaPlaybackDevice{dev}
	} else {
		var dev *alsa.CaptureDevice
		if dev, err = alsa.NewCaptureDevice(device, int(channels), alsa.FormatS16LE, int(samplerate), alsa.BufferParams{}); err != nil {
			return
		}
		alsaInterface = &AlsaCaptureDevice{dev}
	}
	ab.AlsaInterface = alsaInterface
	return
}

func (ab *AlsaBackend) Read(buf []byte) (int, error) {
	return ab.AlsaInterface.Read(buf)
}

func (ab *AlsaBackend) Write(buf []byte) (int, error) {
	return ab.AlsaInterface.Write(buf)
}
