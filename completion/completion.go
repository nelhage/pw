package completion

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var completionLog = log.New(os.Stderr, "completion: ", log.LstdFlags)

type CommandLine []string
type Completer func(CommandLine) ([]string, bool)

func CompleteIfRequested(flags *flag.FlagSet, completer Completer) {
	if len(os.Args) <= 1 || os.Args[1] != "-do-completion" {
		return
	}
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

	cl := parseLineForCompletion(line, int(point))

	for _, word := range getCompletions(cl, flags, completer) {
		fmt.Println(word)
	}
	os.Exit(0)
}
func (c CommandLine) CurrentWord() string {
	return c[len(c)-1]
}

func parseLineForCompletion(line string, point int) CommandLine {
	var cl CommandLine
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
					cl = append(cl, string(word))
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

	return append(cl, string(word))
}

type boolFlag interface {
	flag.Value
	IsBoolFlag() bool
}

func completeFlags(cl CommandLine, flags *flag.FlagSet) (completions []string, rest CommandLine) {
	if len(cl) == 0 {
		return nil, cl
	}
	cl = cl[1:]
	var inFlag string
	for len(cl) > 1 {
		w := cl[0]
		if inFlag != "" {
			inFlag = ""
		} else if len(w) > 1 && w[0] == '-' && w != "--" {
			if !strings.Contains(w, "=") {
				var i int
				for i = 0; i < len(w) && w[i] == '-'; i++ {
				}
				inFlag = w[i:]
			}
			if flag := flags.Lookup(inFlag); flag != nil {
				if bf, ok := flag.Value.(boolFlag); ok && bf.IsBoolFlag() {
					inFlag = ""
				}
			}
		} else {
			if w == "--" {
				cl = cl[1:]
			}
			return nil, cl
		}
		cl = cl[1:]
	}

	if inFlag != "" {
		// Complete a flag value. No-op for now.
		return []string{}, nil
	} else if len(cl[0]) > 0 && cl[0][0] == '-' {
		// complete a flag name
		prefix := strings.TrimLeft(cl[0], "-")
		flags.VisitAll(func(f *flag.Flag) {
			if strings.HasPrefix(f.Name, prefix) {
				completions = append(completions, "-"+f.Name)
			}
		})
		return completions, nil
	}

	if cl[0] == "" {
		flags.VisitAll(func(f *flag.Flag) {
			completions = append(completions, "-"+f.Name)
		})
	}
	return completions, cl
}

func getCompletions(cl CommandLine, flags *flag.FlagSet, completer Completer) []string {
	completions, rest := completeFlags(cl, flags)
	if rest != nil {
		if extra, ok := completer(rest); ok {
			completions = append(completions, extra...)
		}
	}

	return completions
}
