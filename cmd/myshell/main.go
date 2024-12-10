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
	fmt.Fprint(os.Stdout, "$ ")
	for scanner.Scan() {
		inp := scanner.Text()
		out := ""
		tokens := strings.Split(inp, " ")
		// fmt.Println(strconv.Quote(inp))
		switch tokens[0] {
		case constants.EXIT:
			code := 0
			if len(tokens) > 1 {
				code, _ = strconv.Atoi(tokens[1])
			}
			os.Exit(code)
		case constants.ECHO:
			remainBuilder := strings.Builder{}
			remainBuilder.WriteString(strings.Join(tokens[1:], " "))
			hasStart := strings.HasPrefix(remainBuilder.String(), "'")
			hasTerminate := strings.HasSuffix(remainBuilder.String(), "'") && !strings.HasSuffix(remainBuilder.String(), "\\'")
			switch {
			// case hasStart && !hasTerminate:
			// 	remainScanner := bufio.NewScanner(os.Stdin)
			// 	remainScanner.Split(bufio.ScanRunes)
			// 	for remainScanner.Scan() {
			// 		curBytes := scanner.Bytes()
			// 		remainBuilder.Write(curBytes)
			// 		hasTerminate := curBytes[len(curBytes)-1] == '\'' && curBytes[len(curBytes)]
			// 		if hasTerminate {
			// 			break
			// 		}
			// 	}
			// 	fallthrough
			case hasStart && hasTerminate:
				out = strings.Trim(remainBuilder.String(), "'")
			default:
				out = remainBuilder.String()
			}
		case constants.TYPE:
			path, isExists := constants.GetCommand(tokens[1])
			if isExists {
				out = fmt.Sprintf("%v is %v", tokens[1], path)
			} else {
				out = fmt.Sprintf("%v: not found", tokens[1])
			}
		case constants.CD:
			tokens[1] = strings.Replace(tokens[1], "~", os.Getenv("HOME"), 1)
			var path string
			if filepath.IsAbs(tokens[1]) {
				path = tokens[1]
			} else {
				path = filepath.Join(curWorkingDir, tokens[1])
			}
			// fmt.Println(tmp)
			if _, err := os.Stat(path); err != nil {
				out = fmt.Sprintf("cd: %v: No such file or directory", tokens[1])
			} else {
				curWorkingDir = path
			}
		case constants.PWD:
			out = curWorkingDir
		default:
			// fmt.Println(constants.MapCommand2Path)
			program, isExisted := constants.GetCommand(tokens[0])
			if !isExisted {
				out = fmt.Sprintf("%v: command not found", tokens[0])
			} else {
				var args []string
				if len(tokens) > 1 {
					args = tokens[1:]
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
