package completion

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var completionLog = log.New(os.Stderr, "completion: ", log.LstdFlags)

type Completer func(CompletionState) []string

func CompleteIfRequested(completer Completer) {
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

	options, err := doCompletion(line, int(point), completer)
	if err != nil {
		completionLog.Println("Error in completion: ", err.Error())
		os.Exit(1)
	}
	for _, word := range options {
		fmt.Println(word)
	}
	os.Exit(0)
}

type CompletionState struct {
	Line  string
	Point int
	Words []string
}

func (c CompletionState) CurrentWord() string {
	return c.Words[len(c.Words)-1]
}

func parseLineForCompletion(line string, point int) CompletionState {
	state := CompletionState{
		Line:  line,
		Point: point,
	}
	var quote rune = 0
	var backslash bool = false
	var word []rune
	for _, char := range line[:point] {
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
					state.Words = append(state.Words, string(word))
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

	state.Words = append(state.Words, string(word))

	return state
}

func doCompletion(line string, point int, completer Completer) ([]string, error) {
	state := parseLineForCompletion(line, point)
	return completer(state), nil
}
