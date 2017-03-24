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

func (fb *FileBackend) Read(buf []byte) (n int, err error) {
	if n, err = fb.File.Read(buf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	return
}

func (fb *FileBackend) Write(buf []byte) (n int, err error) {
	if n, err = fb.File.Write(buf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	return
}
