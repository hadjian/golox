package main

import (
	"log"
	"strconv"
)

func err(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	log.Println("[line " + strconv.Itoa(line) + "] Error" + where + ": " + message)
	hadError = true
}
