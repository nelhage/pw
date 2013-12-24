package main

import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var getCommand command = command{
	command: "get",
	action:  doGetPassword,
}

func doGetPassword(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: %s get [PASSWORD]", os.Args[0])
	}
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
}

func doEditPassword(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: %s edit [PASSWORD]", os.Args[0])
	}

	decrypted, err := config.ReadPassword(args[0])
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Reading password: %s\n", err)
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
}

func doRmPassword(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: %s rm PASSWORD", os.Args[0])
	}

	if err := config.RemovePassword(args[0]); err != nil {
		return err
	}

	log.Printf("Removed password `%s'\n", args[0])
	return nil
}

var addCommand command = command{
	command: "add",
	action:  doAddPassword,
}

func doAddPassword(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Usage: %s add PASSWORD", os.Args[0])
	}

	fmt.Printf("Contents for password `%s': ", args[0])
	plaintext, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println()

	if err = config.WritePassword(args[0], string(plaintext)); err != nil {
		return err
	}

	log.Printf("Saved password `%s'\n", args[0])
	return nil
}

var lsCommand command = command{
	command: "ls",
	action:  doLsPasswords,
}

func doLsPasswords(args []string) error {
	passwords, err := config.ListPasswords()
	if err != nil {
		return err
	}

	for _, pass := range passwords {
		fmt.Printf("%s\n", pass)
	}

	return nil
}
