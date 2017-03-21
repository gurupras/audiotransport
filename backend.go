package audiotransport

type BackendInterface interface {
	Init(name, device string, samplerate, channels, isPlayback int32) error
	Read(b []byte, len int32) int32
	Write(b []byte, len int32) int32
	GetBufferSize() int32
	GetLatency() int32
}

type Backend struct {
	Name       string
	Device     string
	HandleIdx  int32
	SampleRate int32
	Channels   int32
}

func (b *Backend) Init(name, device string, samplerate, channels int32) {
	b.Name = name
	b.Device = device
	b.SampleRate = samplerate
	b.Channels = channels
}

func (b *Backend) GetBufferSize() int32 {
	return 512
}

func (b *Backend) GetLatency() int32 {
	return -1
}
