package audiotransport

import "github.com/mesilliac/pulse-simple"

type PulseBackend struct {
	*Backend
	*pulse.Stream
}

func (pb *PulseBackend) Init(name, device string, samplerate, channels uint32, isPlayback bool) (err error) {
	if pb.Backend == nil {
		pb.Backend = &Backend{}
	}
	pb.Backend.Init(name, device, samplerate, channels)

	var dir pulse.StreamDirection
	if isPlayback {
		dir = pulse.STREAM_PLAYBACK
	} else {
		dir = pulse.STREAM_RECORD
	}

	spec := &pulse.SampleSpec{pulse.SAMPLE_S16LE, samplerate, uint8(channels)}
	pb.Stream, err = pulse.NewStream("", name, dir, device, "", spec, nil, nil)
	return
}

func (pb *PulseBackend) Read(buf []byte, len uint32) (int, error) {
	return pb.Stream.Read(buf)
}

func (pb *PulseBackend) Write(buf []byte, len uint32) (int, error) {
	return pb.Stream.Write(buf)
}

func (pb *PulseBackend) GetLatency() (int64, error) {
	lat, err := pb.Stream.Latency()
	return int64(lat), err
}
