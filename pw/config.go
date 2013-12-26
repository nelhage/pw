package pw

import (
	"flag"
	"fmt"
	"github.com/nelhage/go.config"
	"os"
)

type Config struct {
	GPGKey      string
	RootDir     string
	CopyCommand string
}

var theConfig Config

func init() {
	flag.StringVar(&theConfig.GPGKey, "gpgkey", "",
		"The GPG key to encrypt passwords to")
	flag.StringVar(&theConfig.RootDir, "root", os.ExpandEnv("${HOME}/pw"),
		"The root directory for the password store")
	flag.StringVar(&theConfig.CopyCommand, "cmd.copy", "xclip -i",
		"A command that accepts input on STDIN and copies it to the clipboard")
}

func LoadConfig() *Config {
	if err := config.LoadConfig(flag.CommandLine, "pw"); err != nil {
		panic(fmt.Sprintf("Loading config: %s", err))
	}
	return &theConfig
}

func (config *Config) recipients() []string {
	if config.GPGKey != "" {
		return []string{config.GPGKey}
	}
	return nil
}
