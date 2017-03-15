package audiotransport

func GetBufferSize(samplerate int32) int32 {
	return samplerate / 4
}
