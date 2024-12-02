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
		if tokens[0] == constants.EXIT {
			code := 0
			if len(tokens) > 1 {
				code, _ = strconv.Atoi(tokens[1])
			}
			os.Exit(code)
		} else {
			str := fmt.Sprintf("%v: command not found\n", strings.TrimRight(inp, "\n"))
			fmt.Fprint(os.Stdout, str)
		}
	}
}
