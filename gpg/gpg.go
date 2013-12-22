package gpg

import (
	"bytes"
	"os"
	"os/exec"
)

func runGPGPipe(args []string, in []byte) ([]byte, error) {
	cmd := exec.Command("gpg", args...)
	var out bytes.Buffer
	cmd.Stdin = bytes.NewBuffer(in)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
