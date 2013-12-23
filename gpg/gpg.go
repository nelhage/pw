package gpg

import (
	"bytes"
	"os"
	"os/exec"
)

func runGPGPipe(args []string, in string) (string, error) {
	cmd := exec.Command("gpg", args...)
	var out bytes.Buffer
	cmd.Stdin = bytes.NewBuffer([]byte(in))
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return string(out.Bytes()), nil
}

func GPGEncrypt(keys []string, in string) (string, error) {
	args := []string{"--encrypt", "--armor"}
	for _, keyid := range keys {
		args = append(args, "-r", keyid)
	}
	return runGPGPipe(args, in)
}

func GPGDecrypt(in string) (string, error) {
	return runGPGPipe([]string{"--decrypt"}, in)
}
