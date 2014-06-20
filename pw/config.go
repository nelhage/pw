package pw

import (
	"flag"
	"github.com/nelhage/go.cli/config"
	"os"
	"runtime"
)

type Config struct {
	GPGKey          string
	RootDir         string
	CopyCommand     string
	GenerateCommand string
}

var theConfig Config

func init() {
	flag.StringVar(&theConfig.GPGKey, "gpgkey", "",
		"The GPG key to encrypt passwords to")
	flag.StringVar(&theConfig.RootDir, "root", os.ExpandEnv("${HOME}/pw"),
		"The root directory for the password store")
	var copyCmd string
	switch runtime.GOOS {
	case "linux":
		copyCmd = "xclip -i"
	case "darwin":
		copyCmd = "pbcopy"
	}
	flag.StringVar(&theConfig.CopyCommand, "cmd.copy", copyCmd,
		"A command that accepts input on STDIN and copies it to the clipboard")
	flag.StringVar(&theConfig.GenerateCommand, "cmd.generate", "pwgen 12 1 -s",
		"A command to generate new passwords")
}

func LoadConfig() (*Config, error) {
	if err := config.LoadConfig(flag.CommandLine, "pw"); err != nil {
		return nil, err
	}
	return &theConfig, nil
}

func (config *Config) recipients() []string {
	if config.GPGKey != "" {
		return []string{config.GPGKey}
	}
	return nil
}
