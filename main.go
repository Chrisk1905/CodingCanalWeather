package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("starting REPL..")
	for {
		fmt.Print("weather > ")
		if scanner.Scan() {
			text := scanner.Text()
			split_text := strings.Split(text, " ")
			args := split_text[1:]
			fmt.Printf("%s", args)
		}
	}
}
