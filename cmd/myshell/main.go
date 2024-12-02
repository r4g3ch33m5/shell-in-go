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
	entryUri := filepath.Join(prefix, entry.Name())
	fmt.Println(strings.Repeat("\t", level), entryUri)
	if entry.IsDir() {
		childEntries, _ := os.ReadDir(entry.Name())
		for _, entry := range childEntries {
			readEntry(entry, entryUri, level+1)
		}
	}
	if constants.MapCommand2Path[entry.Name()] != constants.BUILTIN {
		constants.MapCommand2Path[entry.Name()] = entryUri
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

				program, isExisted := constants.MapCommand2Path[tokens[0]]
				if !isExisted {
					pathEnv := os.Getenv("PATH")
					paths := strings.Split(pathEnv, ":")
					for _, path := range paths {
						fileEntries, _ := os.ReadDir(filepath.Dir(path))
						for _, entry := range fileEntries {
							readEntry(entry, path, 0)
						}
					}
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
