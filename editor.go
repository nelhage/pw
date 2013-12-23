package pw

import (
	"os"
	"os/exec"
	"strings"
)

func (config *Config) RunEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vi"
	}

	quotedPath := "'" + strings.Replace(path, "'", `'\''`, -1) + "'"

	cmd := exec.Command("sh", "-c", editor+" "+quotedPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
