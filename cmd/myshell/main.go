package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/constants"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func isSpace(r byte) bool {
	// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
	switch r {
	case ' ', '\t', '\n', '\v', '\f', '\r':
		return true
	case '\u0085', '\u00A0':
		return true
	}
	return false
}

type Opt struct {
	KeepQuote bool
}

func retrieveArgs(scanner *bufio.Scanner, opts ...Opt) strings.Builder {
	// hasSlash := false
	hasQuote := false
	hasDQuote := false
	buffer := strings.Builder{}
	opt := Opt{}
	if len(opts) > 0 {
		opt = opts[0]
	}

bufferScan:
	for {
		switch scanner.Bytes()[0] {
		case '\r', '\n':
			if hasQuote || hasDQuote {
				buffer.Write(scanner.Bytes())
			} else {
				break bufferScan
			}
		case '\'':
			if hasDQuote {
				buffer.Write(scanner.Bytes())
			} else {
				if opt.KeepQuote {
					buffer.Write(scanner.Bytes())
				}
				hasQuote = !hasQuote
			}
		case '"':
			if hasQuote {
				buffer.Write(scanner.Bytes())
			} else {
				if opt.KeepQuote {
					buffer.Write(scanner.Bytes())
				}
				hasDQuote = !hasDQuote
			}
		case ' ', '\t':
			if hasDQuote || hasQuote {
				buffer.Write(scanner.Bytes())
			} else if len(buffer.String()) != 0 {
				buffer.Write(scanner.Bytes())
				// skip trailing
				for scanner.Scan() {
					if !(scanner.Bytes()[0] == ' ' || scanner.Bytes()[0] == '\t') {
						continue bufferScan
					}
				}
			}
		default:
			buffer.Write(scanner.Bytes())
		}
		scanner.Scan()
	}
	return buffer
}

func main() {
	// Uncomment this block to pass the first stage
	// inpChan := make(chan string)
	// outputChan := shell.NewTokenizer(inpChan)
	// Wait for user input
	curWorkingDir, _ := filepath.Abs(".")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanBytes)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		for scanner.Bytes()[0] == '\r' || scanner.Bytes()[0] == '\n' {
			fmt.Fprint(os.Stdout, "$ ")
			scanner.Scan()
			continue
		}
		cmdBuffer := strings.Builder{}
		for !isSpace(scanner.Bytes()[0]) {
			cmdBuffer.Write(scanner.Bytes())
			scanner.Scan()
		}
		cmd := cmdBuffer.String()
		var out string
		switch cmd {
		case constants.EXIT:
			code := 0
			buffer := retrieveArgs(scanner)
			tokens := buffer.String()
			if len(tokens) != 0 {
				code, _ = strconv.Atoi(tokens)
			}
			os.Exit(code)
		case constants.ECHO:
			buffer := retrieveArgs(scanner)
			out = buffer.String()
		case constants.TYPE:
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			tokens := scanner.Text()
			path, isExists := constants.GetCommand(tokens)
			if isExists {
				out = fmt.Sprintf("%v is %v", tokens, path)
			} else {
				out = fmt.Sprintf("%v: not found", tokens)
			}
		case constants.CD:
			buffer := retrieveArgs(scanner)
			tokens := buffer.String()
			tokens = strings.Replace(tokens, "~", os.Getenv("HOME"), 1)
			var path string
			if filepath.IsAbs(tokens) {
				path = tokens
			} else {
				path = filepath.Join(curWorkingDir, tokens)
			}
			// fmt.Println(tmp)
			if _, err := os.Stat(path); err != nil {
				out = fmt.Sprintf("cd: %v: No such file or directory", tokens)
			} else {
				curWorkingDir = path
			}
		case constants.PWD:
			out = curWorkingDir
		default:
			// fmt.Println(constants.MapCommand2Path)
			// scanner := bufio.NewScanner(os.Stdin)
			// scanner.Scan()
			program, isExisted := constants.GetCommand(cmd)
			// retrieve args
			buffer := retrieveArgs(scanner, Opt{KeepQuote: true})
			if !isExisted {
				out = fmt.Sprintf("%v: command not found", cmd)
			} else {
				tokens := buffer.String()
				quoted := false
				args := strings.FieldsFunc(tokens, func(r rune) bool {
					if r == '"' {
						quoted = !quoted
					}
					return !quoted && r == ' '
				})
				// fmt.Println(strconv.Quote(tokens))
				for idx, arg := range args {
					args[idx], _ = strconv.Unquote(arg)
				}
				command := exec.Command(program, args...)
				command.Stderr = os.Stderr
				command.Stdout = os.Stdout
				err := command.Run()
				if err != nil {
					fmt.Println(command.String())
				}
				// out = strings.TrimSuffix(string(output), "\n")
			}
		}
		if len(out) != 0 {
			out += "\n"
		}
		// fmt.Println("quote out: ", strconv.Quote(out))
		fmt.Fprint(os.Stdout, out, "$ ")
	}
}
