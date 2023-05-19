package main

import (
	"log"
	"strconv"
)

func errLine(line int, message string) {
	report(line, "", message)
}

func errToken(token Token, message string) {
	if token.tType == EOF {
		report(token.line, " at end", message)
	} else {
		report(token.line, " at '"+token.lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	log.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	hadError = true
}
