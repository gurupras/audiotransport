package audiotransport

type ApiType int

const (
	ALSA_API  ApiType = iota
	PULSE_API ApiType = iota
)

func (apiType ApiType) ApiString() string {
	switch apiType {
	case ALSA_API:
		return "ALSA_API"
	case PULSE_API:
		return "PULSE_API"
	}
	return ""
}

func (apiType ApiType) GetBufferSize(samplerate int32, channels int32) int32 {
	switch apiType {
	case ALSA_API:
		return 128 * (16 / 8) * 2
	case PULSE_API:
		return 512
	}
	return -1
}
