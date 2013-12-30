package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"github.com/nelhage/pw/pw"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var getCommand command = command{
	command: "get",
	action:  doGetPassword,
	usage:   "PASSWORD",
	minArgs: 1,
}

func doGetPassword(args []string) error {
	decrypted, err := config.ReadPassword(args[0])
	if err != nil {
		return fmt.Errorf("Unable to read password: %s\n", err.Error())
	}
	log.Printf("Contents of `%s':\n", args[0])
	fmt.Printf("%s", decrypted)
	return nil
}

var editCommand command = command{
	command: "edit",
	action:  doEditPassword,
	minArgs: 1,
	usage:   "PASSWORD",
}

func doEditPassword(args []string) error {
	decrypted, err := config.ReadPassword(args[0])
	if err != nil {
		if _, ok := err.(*pw.NoSuchPassword); !ok {
			decrypted = ""
		} else {
			return err
		}
	}

	f, err := ioutil.TempFile("", "pw-")
	if err != nil {
		return fmt.Errorf("Creating temp file: %s", err.Error())
	}
	defer func() {
		// TODO: Secure delete?
		f.Close()
		os.Remove(f.Name())
	}()
	io.WriteString(f, decrypted)

	if err = config.RunEditor(f.Name()); err != nil {
		return fmt.Errorf("Running editor: %s\n", err)
	}

	if _, err = f.Seek(0, 0); err != nil {
		return err
	}

	newPlaintext, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = config.WritePassword(args[0], string(newPlaintext)); err != nil {
		return err
	}

	return nil
}

var rmCommand command = command{
	command: "rm",
	action:  doRmPassword,
	minArgs: 1,
	usage:   "PASSWORD",
}

func doRmPassword(args []string) error {
	if err := config.RemovePassword(args[0]); err != nil {
		return err
	}

	log.Printf("Removed password `%s'\n", args[0])
	return nil
}

var addCommand command = command{
	command: "add",
	action:  doAddPassword,
	minArgs: 1,
	usage:   "PASSWORD [FILE]",
}

func doAddPassword(args []string) error {
	var plaintext []byte
	var err error

	if len(args) == 1 && terminal.IsTerminal(0) {
		fmt.Printf("Contents for password `%s': ", args[0])
		plaintext, err = terminal.ReadPassword(0)
		if err != nil {
			return err
		}
		fmt.Println()
	} else {
		var f io.Reader
		if len(args) > 1 {
			file, err := os.Open(args[1])
			if err != nil {
				return err
			}
			defer file.Close()
			f = file
		} else {
			f = os.Stdin
		}
		if plaintext, err = ioutil.ReadAll(f); err != nil {
			return err
		}
	}

	if err = config.WritePassword(args[0], string(plaintext)); err != nil {
		return err
	}

	log.Printf("Saved password `%s'\n", args[0])
	return nil
}

var lsCommand command = command{
	command: "ls",
	action:  doLsPasswords,
	usage:   "[NEEDLE]",
}

func doLsPasswords(args []string) error {
	var filter func(string) bool
	if len(args) == 0 {
		filter = func(_ string) bool { return true }
	} else {
		filter = func(pw string) bool { return strings.Contains(pw, args[0]) }
	}
	passwords, err := config.ListPasswords()
	if err != nil {
		return err
	}

	for _, pass := range passwords {
		if filter(pass) {
			fmt.Printf("%s\n", pass)
		}
	}

	return nil
}

var copyCommand command = command{
	command: "copy",
	action:  doCopyPassword,
	minArgs: 1,
	usage:   "PASSWORD",
}

func doCopyPassword(args []string) error {
	plaintext, err := config.ReadPassword(args[0])
	if err != nil {
		return err
	}

	if strings.Count(plaintext, "\n") == 1 {
		plaintext = strings.TrimSuffix(plaintext, "\n")
	}

	if err = config.CopyText(plaintext); err != nil {
		return err
	}

	log.Printf("Copied password `%s' to clipboard.\n", args[0])
	return nil
}

var newCommand command = command{
	command: "new",
	action:  doNewPassword,
	minArgs: 1,
	usage:   "PASSWORD",
}

func doNewPassword(args []string) error {
	password, err := config.GeneratePassword()

	if err != nil {
		return err
	}

	if err := config.WritePassword(args[0], password); err != nil {
		return err
	}

	if err := config.CopyText(password); err != nil {
		return err
	}

	log.Printf("Generated password `%s' and copied to the clipboard.\n", args[0])
	return nil
}
