package audiotransport

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gurupras/gocommons"
)

type FileHeader struct {
	ID     [4]byte
	Length int
	Type   [4]byte
}

func ParseFileHeader(in io.Reader) (fh *FileHeader, err error) {
	var n int
	b := make([]byte, 12)
	if n, err = in.Read(b); err != nil {
		return
	} else if n < 12 {
		err = errors.New(fmt.Sprintf("Expected length of at least 12; got %v", n))
		return
	}
	fh = &FileHeader{}
	copy(fh.ID[:], b[:4])
	fh.Length = int(binary.LittleEndian.Uint32(b[4:8]))
	copy(fh.Type[:], b[8:12])

	return
}

type ChunkHeader struct {
	ID     [4]byte
	Length int
}

func ParseChunkHeader(in io.Reader) (ch *ChunkHeader, err error) {
	var n int
	b := make([]byte, 8)
	if n, err = in.Read(b); err != nil {
		return
	} else if n < 8 {
		err = errors.New(fmt.Sprintf("Expected length of at least 8; got %v", n))
		return
	}
	ch = &ChunkHeader{}
	copy(ch.ID[:], b[:4])
	ch.Length = int(binary.LittleEndian.Uint32(b[4:8]))
	return
}

type FormatHeader struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    int
	ByteRate      int
	BlockAlign    uint16
	BitsPerSample uint16
}

func ParseFormatHeader(in io.Reader) (fh *FormatHeader, err error) {
	var n int
	b := make([]byte, 16)
	if n, err = in.Read(b); err != nil {
		return
	} else if n < 16 {
		err = errors.New(fmt.Sprintf("Expected length of at least 16; got %v", n))
		return
	}

	fh = &FormatHeader{}
	fh.AudioFormat = binary.LittleEndian.Uint16(b[:2])
	fh.NumChannels = binary.LittleEndian.Uint16(b[2:4])
	fh.SampleRate = int(binary.LittleEndian.Uint32(b[4:8]))
	fh.ByteRate = int(binary.LittleEndian.Uint32(b[8:12]))
	fh.BlockAlign = binary.LittleEndian.Uint16(b[12:14])
	fh.BitsPerSample = binary.LittleEndian.Uint16(b[14:16])
	return
}

type WavHeader struct {
	*FileHeader
	*ChunkHeader
	*FormatHeader
}

func ParseWavHeaders(in io.Reader) (wh *WavHeader, err error) {
	wh = &WavHeader{}
	if wh.FileHeader, err = ParseFileHeader(in); err != nil {
		return nil, err
	}
	if wh.ChunkHeader, err = ParseChunkHeader(in); err != nil {
		return nil, err
	}
	if wh.FormatHeader, err = ParseFormatHeader(in); err != nil {
		return nil, err
	}
	return
}

func ParseWavFile(path string) (wh *WavHeader, data []byte, err error) {
	var file *gocommons.File
	var dataHeader *ChunkHeader

	if file, err = gocommons.Open(path, os.O_RDONLY, gocommons.GZ_FALSE); err != nil {
		return
	}

	if wh, err = ParseWavHeaders(file.File); err != nil {
		return
	}

	if dataHeader, err = ParseChunkHeader(file.File); err != nil {
		return
	}

	data = make([]byte, dataHeader.Length)
	_, err = file.File.Read(data)
	return
}
