package main

import (
	"flag"
	"fmt"
	"github.com/nelhage/pw/pw"
	"log"
	"os"
	"strings"
)

var config *pw.Config = pw.LoadConfig()

type command struct {
	command string
	action  func([]string) error
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

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] COMMAND [ARGS]\n", os.Args[0])
		fmt.Printf(" Known commands: %s\n", strings.Join(knownCommands(), ", "))
		flag.PrintDefaults()
	}

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	for _, cmd := range commands {
		if args[0] == cmd.command {
			if err := cmd.action(args[1:]); err != nil {
				log.Fatalln(err.Error())
			}
			return
		}
	}

	log.Printf("Unknown command: %s\n", args[0])
	flag.Usage()
	os.Exit(1)
}
