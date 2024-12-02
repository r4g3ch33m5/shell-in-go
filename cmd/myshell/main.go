package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/constants"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func readEntry(entry fs.DirEntry, prefix string, level int) {
	curEntry := filepath.Join(prefix, entry.Name())
	// fmt.Println(strings.Repeat("| ", level), entry.Name())
	if constants.MapCommand2Path[entry.Name()] != constants.BUILTIN {
		fmt.Println(curEntry)
		constants.MapCommand2Path[entry.Name()] = curEntry
	}
	if entry.IsDir() {
		childEntries, err := os.ReadDir(curEntry)
		if err != nil {
			fmt.Println(strings.Repeat("| ", level+1), entry.Name(), err)
		}
		for _, entry := range childEntries {
			readEntry(entry, curEntry, level+1)
		}
	}
}

func main() {
	// Uncomment this block to pass the first stage
	// inpChan := make(chan string)
	// outputChan := shell.NewTokenizer(inpChan)
	// Wait for user input
	out := ""

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Fprint(os.Stdout, "$ ")
		for scanner.Scan() {
			inp := scanner.Text()
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
				out = strings.Join(tokens[1:], " ")
			case constants.TYPE:
				path, isExists := constants.MapCommand2Path[tokens[1]]
				if isExists {
					out = fmt.Sprintf("%v is %v", tokens[1], path)
				} else {
					out = fmt.Sprintf("%v: not found", tokens[1])
				}
			default:
				pathEnv := os.Getenv("PATH")
				paths := strings.Split(pathEnv, ":")
				for _, path := range paths {
					fileEntries, _ := os.ReadDir(path)
					for _, entry := range fileEntries {
						readEntry(entry, path, 0)
					}
				}
				fmt.Println(constants.MapCommand2Path)
				program, isExisted := constants.MapCommand2Path[tokens[0]]
				if !isExisted {
					out = fmt.Sprintf("%v: command not found", tokens[0])
				} else {
					var args []string
					if len(tokens) > 1 {
						args = tokens[1:]
					}
					exec.Command(program, args...)
				}
			}
			fmt.Fprint(os.Stdout, out, "\n$ ")
		}
	}
}
