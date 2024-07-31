package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mactavishz/re-nfa/pkg/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./re <regex> [<string>]")
		fmt.Println("If <string> is not provided, the program will read from stdin.")
		os.Exit(1)
	}

	regex := os.Args[1]
	var input string

	if len(os.Args) == 3 {
		// Input string provided as an argument
		input = os.Args[2]
	} else {
		// Read input from stdin
		reader := bufio.NewReader(os.Stdin)
		inputBytes, err := reader.ReadBytes('\n')
		if err != nil && err.Error() != "EOF" {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(string(inputBytes))
	}

	parser := utils.NewParser(regex)
	nfa, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error parsing regex '%s': %v\n", regex, err)
		os.Exit(1)
	}

	matched := nfa.Match(input)
	if matched {
		fmt.Printf("%s\n", input)
	}
}
