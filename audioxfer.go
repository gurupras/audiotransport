package audiotransport

type ApiType int

const (
	ALSA_API  ApiType = iota
	PULSE_API ApiType = iota
	FILE_API  ApiType = iota
)

func (apiType ApiType) ApiString() string {
	switch apiType {
	case ALSA_API:
		return "ALSA_API"
	case PULSE_API:
		return "PULSE_API"
	case FILE_API:
		return "FILE_API"
	}
	return ""
}
