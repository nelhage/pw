package pw

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	GPGKey  string
	RootDir string
}

var theConfig Config

func init() {
	flag.StringVar(&theConfig.GPGKey, "gpgkey", "", "The GPG key to encrypt passwords to")
	flag.StringVar(&theConfig.RootDir, "root", os.ExpandEnv("${HOME}/pw"), "The root directory for the password store")
}

func LoadConfig() *Config {
	if err := theConfig.loadConfig(); err != nil {
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

func (config *Config) loadConfig() error {
	path := os.ExpandEnv("${HOME}/.pw")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return config.parseConfig(f)
}

func (config *Config) parseConfig(f io.Reader) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		bits := strings.SplitN(line, "=", 2)
		if len(bits) != 2 {
			return fmt.Errorf("Illegal config line: `%s'", line)
		}

		key := strings.TrimSpace(bits[0])
		value := strings.TrimSpace(bits[1])

		if flag := flag.CommandLine.Lookup(key); flag == nil {
			return fmt.Errorf("Unknown option `%s'", bits[0])
		}

		if err := flag.CommandLine.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}
