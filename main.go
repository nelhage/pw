package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nelhage/go.cli/completion"
	"github.com/nelhage/pw/pw"
)

var config *pw.Config

type command struct {
	command string
	action  func([]string) error
	usage   string
	minArgs int
}

var commands []command = []command{
	getCommand,
	editCommand,
	rmCommand,
	addCommand,
	lsCommand,
	copyCommand,
	newCommand,
}

func knownCommands() []string {
	var out []string
	for _, cmd := range commands {
		out = append(out, cmd.command)
	}
	return out
}

func runCommand(cmd *command, args []string) error {
	if len(args) < cmd.minArgs {
		var pad string
		if cmd.usage == "" {
			pad = ""
		} else {
			pad = " "
		}
		return fmt.Errorf("Usage: %s %s%s%s",
			os.Args[0],
			cmd.command,
			pad,
			cmd.usage)
	}
	return cmd.action(args)
}

func complete(cl completion.CommandLine) (completions []string) {
	switch len(cl) {
	case 1:
		return completion.SetCompleter(knownCommands()).Complete(cl)
	case 2:
		passwords, err := config.ListPasswords()
		if err != nil {
			return nil
		}
		return completion.SetCompleter(passwords).Complete(cl)
	}
	return completions
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] COMMAND [ARGS]\n", os.Args[0])
		fmt.Printf(" Known commands: %s\n", strings.Join(knownCommands(), ", "))
		flag.PrintDefaults()
	}
	var err error
	config, err = pw.LoadConfig()
	if err != nil {
		log.Fatalf("Loading config: %s", err)
	}

	completion.CompleteIfRequested(
		completion.CompleterWithFlags(flag.CommandLine,
			completion.FunctionCompleter(complete)))

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if args[0] == cmd.command {
			if err := runCommand(&cmd, args[1:]); err != nil {
				log.Fatalln(err.Error())
			}
			return
		}
	}

	log.Printf("Unknown command: %s\n", args[0])
	flag.Usage()
	os.Exit(1)
}
