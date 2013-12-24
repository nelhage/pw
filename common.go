package pw

import (
	"fmt"
	"github.com/nelhage/pw/gpg"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	if !strings.HasSuffix(plaintext, "\n") {
		plaintext = plaintext + "\n"
	}

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

func (config *Config) RemovePassword(pw string) error {
	fullPath := config.ResolvePath(pw)

	if err := os.Remove(fullPath); err != nil {
		p := err.(*os.PathError)
		if os.IsNotExist(p.Err) {
			return fmt.Errorf("No such password: %s", err)
		}
		return err
	}
	return nil
}

func (config *Config) ListPasswords() ([]string, error) {
	var passwords []string
	err := filepath.Walk(config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() && strings.HasSuffix(path, ".gpg") {
			pass := strings.TrimPrefix(path, config.RootDir)
			pass = strings.TrimPrefix(pass, "/")
			pass = strings.TrimSuffix(pass, ".gpg")
			passwords = append(passwords, pass)
		}

		return nil
	})
	return passwords, err
}
