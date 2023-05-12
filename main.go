package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Println(err)
		}
	}
}

func runFile(f string) error {
	data, e := ioutil.ReadFile(f)
	if e != nil {
		return e
	}
	run(string(data))
	if hadError {
		os.Exit(65)
	}
	return nil
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		scanner.Scan()
		line := scanner.Text()
		if (line == "") || scanner.Err() != nil {
			fmt.Println("Bye!")
			break
		}
		run(line)
		fmt.Printf("\n")
		hadError = false
	}
	return scanner.Err()
}

func run(script string) {
	fmt.Println(script)
}
