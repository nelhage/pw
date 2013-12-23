package pw

import (
	"github.com/nelhage/pw/gpg"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func (config *Config) ResolvePath(pw string) string {
	return path.Join(config.RootDir, pw) + ".gpg"
}

func (config *Config) ReadPassword(pw string) (string, error) {
	fullPath := config.ResolvePath(pw)

	f, err := os.Open(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return gpg.GPGDecrypt(string(data))
}

func (config *Config) WritePassword(pw string, plaintext string) error {
	fullPath := config.ResolvePath(pw)

	encrypted, err := gpg.GPGEncrypt([]string{config.GPGKey}, plaintext)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(path.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.WriteString(f, encrypted); err != nil {
		return err
	}
	return nil
}
