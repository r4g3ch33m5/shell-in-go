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
	for {
		fmt.Fprint(os.Stdout, "$ ")
		inp, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		tokens := strings.Split(inp, " ")
		switch tokens[0] {
		case constants.EXIT:
			code := 0
			if len(tokens) > 1 {
				code, _ = strconv.Atoi(tokens[1])
			}
			os.Exit(code)
		case constants.ECHO:
			fmt.Fprint(os.Stdout, strings.TrimSpace(inp[len(tokens):]))
		default:
			str := fmt.Sprintf("%v: command not found\n", strings.TrimRight(inp, "\n"))
			fmt.Fprint(os.Stdout, str)
		}
	}
}
