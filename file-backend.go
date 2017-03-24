package audiotransport

import (
	"fmt"
	"os"
)

type FileBackend struct {
	*Backend
	*os.File
}

func (fb *FileBackend) Init(name, device string, samplerate, channels uint32, isPlayback bool) (err error) {
	if fb.Backend == nil {
		fb.Backend = &Backend{}
	}
	fb.Backend.Init(name, device, samplerate, channels)

	var flags int
	if isPlayback {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		flags = os.O_RDONLY
	}
	fmt.Printf("FLAGS=%X\n", flags)
	fb.File, err = os.OpenFile(fb.Backend.Name, flags, 0644)
	return
}

func (fb *FileBackend) Read(buf []byte, len uint32) (int, error) {
	n, err := fb.File.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	_ = n
	return 0, err
}

func (fb *FileBackend) Write(buf []byte, len uint32) (int, error) {
	n, err := fb.File.Write(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	_ = n
	return 0, err
}
