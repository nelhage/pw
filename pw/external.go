package pw

import (
	"bytes"
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

func (config *Config) CopyText(text string) error {
	cmd := exec.Command("sh", "-c", config.CopyCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func (config *Config) GeneratePassword() (string, error) {
	var buffer bytes.Buffer
	cmd := exec.Command("sh", "-c", config.GenerateCommand)
	cmd.Stdout = &buffer
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimRight(buffer.String(), "\n"), nil
}
