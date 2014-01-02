package main

import (
	"flag"
	"fmt"
	"github.com/nelhage/pw/completion"
	"github.com/nelhage/pw/pw"
	"log"
	"os"
	"strings"
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

func complete(cl completion.CommandLine) (completions []string, ok bool) {
	switch len(cl) {
	case 1:
		for _, cmd := range commands {
			if strings.HasPrefix(cmd.command, cl.CurrentWord()) {
				completions = append(completions, cmd.command)
			}
		}
		ok = true
	case 2:
		passwords, err := config.ListPasswords()
		if err != nil {
			return nil, false
		}
		for _, pw := range passwords {
			if strings.HasPrefix(pw, cl.CurrentWord()) {
				completions = append(completions, pw)
			}
		}
		ok = true
	}
	return completions, ok
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] COMMAND [ARGS]\n", os.Args[0])
		fmt.Printf(" Known commands: %s\n", strings.Join(knownCommands(), ", "))
		flag.PrintDefaults()
	}
	config = pw.LoadConfig()

	completion.CompleteIfRequested(flag.CommandLine, complete)

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
