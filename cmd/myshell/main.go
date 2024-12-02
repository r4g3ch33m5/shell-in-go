package main

import (
	"bufio"
	"fmt"
	"os"
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
				out = fmt.Sprintf("%v: command not found", tokens[0])
			}
			fmt.Fprint(os.Stdout, out, "\n$ ")
		}
	}
}
