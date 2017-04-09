package audiotransport

import pulse "github.com/mesilliac/pulse-simple"

type PulseSimpleBackend struct {
	*Backend
	*pulse.Stream
}

func (psb *PulseSimpleBackend) Init(name, device string, samplerate, channels uint32, isPlayback bool) (err error) {
	if psb.Backend == nil {
		psb.Backend = &Backend{}
	}
	psb.Backend.Init(name, device, samplerate, channels)

	var dir pulse.StreamDirection
	if isPlayback {
		dir = pulse.STREAM_PLAYBACK
	} else {
		dir = pulse.STREAM_RECORD
	}

	bufAttr := pulse.NewBufferAttr()
	_ = bufAttr
	/*
		bufAttr.Fragsize = psb.Backend.GetBufferSize()
		bufAttr.Maxlength = psb.Backend.GetBufferSize() * 4
		bufAttr.Tlength = bufAttr.Maxlength / 2
		bufAttr.Minreq = psb.Backend.GetBufferSize() / 4
		bufAttr.Prebuf = psb.Backend.GetBufferSize()
	*/
	spec := &pulse.SampleSpec{pulse.SAMPLE_S16LE, samplerate, uint8(channels)}
	psb.Stream, err = pulse.NewStream("", name, dir, device, "", spec, nil, nil)
	return
}

func (psb *PulseSimpleBackend) Read(buf []byte) (int, error) {
	return psb.Stream.Read(buf)
}

func (psb *PulseSimpleBackend) Write(buf []byte) (int, error) {
	return psb.Stream.Write(buf)
}

func (psb *PulseSimpleBackend) GetLatency() (int64, error) {
	lat, err := psb.Stream.Latency()
	return int64(lat), err
}
