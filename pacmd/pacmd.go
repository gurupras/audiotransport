package pacmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gurupras/go-simpleexec"
)

func pacmdList(suffix string) ([]string, error) {
	cmdline := fmt.Sprintf("pacmd list-%v", suffix)
	cmd := simpleexec.ParseCmd(cmdline).Pipe("grep 'name: '").Pipe(`sed -e 's/^[ \t]\+name: \(.*\)$/\1/g'`)
	if cmd == nil {
		return nil, fmt.Errorf("Failed to create pacmd command")
	}
	result := bytes.NewBuffer(nil)
	cmd.Stdout = result

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to start pacmd command: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Failed to wait for pacmd command: %v", err)
	}
	return strings.Split(strings.TrimSpace(result.String()), "\n"), nil
}

func ListSources() ([]string, error) {
	return pacmdList("sources")
}

func ListSinks() ([]string, error) {
	return pacmdList("sinks")
}
