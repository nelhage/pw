package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var completionLog = log.New(os.Stderr, "completion: ", log.LstdFlags)

func maybeComplete() {
	if len(os.Args) <= 1 || os.Args[1] != "-do-completion" {
		return
	}
	// TODO: handle flags
	line := os.Getenv("COMP_LINE")
	pointStr := os.Getenv("COMP_POINT")
	if line == "" || pointStr == "" {
		completionLog.Println("Completion requested, but COMP_LINE and/or COMP_POINT unset.")
		os.Exit(1)
	}

	point, err := strconv.ParseInt(pointStr, 10, 32)
	if err != nil {
		completionLog.Println("Invalid COMP_POINT: ", point)
		os.Exit(1)
	}

	options, err := doCompletion(line, int(point))
	if err != nil {
		completionLog.Println("Error in completion: ", err.Error())
		os.Exit(1)
	}
	for _, word := range options {
		fmt.Println(word)
	}
	os.Exit(0)
}

type completionState struct {
	line  string
	point int
	words []string
	word  int
}

func parseLineForCompletion(line string, point int) completionState {
	state := completionState{
		line:  line,
		point: point,
	}
	var quote rune = 0
	var backslash bool = false
	var word []rune
	for byte, char := range line {
		if state.word == 0 && point <= byte {
			state.word = len(state.words)
		}

		if backslash {
			word = append(word, char)
			backslash = false
			continue
		}
		if char == '\\' {
			word = append(word, char)
			backslash = true
			continue
		}

		switch quote {
		case 0:
			switch char {
			case '\'', '"':
				word = append(word, char)
				quote = char
			case ' ', '\t':
				if word != nil {
					state.words = append(state.words, string(word))
				}
				word = nil
			default:
				word = append(word, char)
			}
		case '\'':
			word = append(word, char)
			if char == '\'' {
				quote = 0
			}
		case '"':
			word = append(word, char)
			if char == '"' {
				quote = 0
			}
		}
	}

	state.words = append(state.words, string(word))

	if point >= len(line) {
		state.word = len(state.words) - 1
	}

	return state
}

func doCompletion(line string, point int) ([]string, error) {
	state := parseLineForCompletion(line, point)
	return complete(state), nil
}

func complete(state completionState) []string {
	var completions []string
	if state.word == 0 {
		return nil
	} else if state.word == 1 {
		for _, cmd := range commands {
			if strings.HasPrefix(cmd.command, state.words[state.word]) {
				completions = append(completions, cmd.command)
			}
		}
	} else if state.word == 2 {
		passwords, err := config.ListPasswords()
		if err != nil {
			return nil
		}
		for _, pw := range passwords {
			if strings.HasPrefix(pw, state.words[state.word]) {
				completions = append(completions, pw)
			}
		}
	}
	return completions
}
