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

func main() {
	// Uncomment this block to pass the first stage
	// inpChan := make(chan string)
	// outputChan := shell.NewTokenizer(inpChan)
	// Wait for user input
	curWorkingDir, _ := filepath.Abs(".")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		cmd := scanner.Text()
		// fmt.Println(strconv.Quote(inp))
		var out string
		switch cmd {
		case constants.EXIT:
			code := 0
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			tokens := scanner.Text()
			if len(tokens) != 0 {
				code, _ = strconv.Atoi(tokens)
			}
			os.Exit(code)
		case constants.ECHO:
			builder := strings.Builder{}
			runeScanner := bufio.NewScanner(os.Stdin)
			runeScanner.Split(bufio.ScanRunes)
			hasQuote := false
			hasSlash := false
		runeScan:
			for runeScanner.Scan() {
				curRune := runeScanner.Bytes()[0]
				switch curRune {
				case '\\':
					hasSlash = !hasSlash
					builder.WriteByte(curRune)
				case '\'':
					if hasSlash {
						builder.WriteByte(curRune)
					}
					hasQuote = !hasQuote
				case '\n':
					if !hasQuote {
						break runeScan
					}
				default:
					builder.WriteByte(curRune)
				}
			}

			out = builder.String()
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
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			tokens := scanner.Text()
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
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			tokens := scanner.Text()
			program, isExisted := constants.GetCommand(cmd)
			if !isExisted {
				out = fmt.Sprintf("%v: command not found", cmd)
			} else {
				var args []string
				if len(tokens) > 1 {
					args = strings.Split(tokens, " ")
				}
				output, _ := exec.Command(program, args...).Output()
				out = strings.TrimSuffix(string(output), "\n")
			}
		}
		if len(out) != 0 {
			out += "\n"
		}
		fmt.Fprint(os.Stdout, out, "$ ")
	}
}
