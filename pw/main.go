package main

import (
	"flag"
	"github.com/nelhage/pw"
	"log"
	"os"
)

var config pw.Config = pw.LoadConfig()

type command struct {
	command string
	action  func([]string) error
}

var commands []command = []command{
	getCommand,
	editCommand,
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("Usage: %s <COMMAND> [ARGS]\n", os.Args[0])
	}

	for _, cmd := range commands {
		if args[0] == cmd.command {
			if err := cmd.action(args[1:]); err != nil {
				log.Fatalln(err.Error())
			}
			return
		}
	}

	log.Fatalf("Unknown command: %s\n", args[0])
}
