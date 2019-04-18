package main

import (
	"os"
	"fmt"
	"log"
)

// Function to get the bots token. It can be either passed as a command-line
// argument, or be read from stdin if no argument is present.
// If neither of those succeed, the function will exit the program.
func token() string {
	tok := ""

	if len(os.Args) < 2 {
		fmt.Scan(&tok)
	} else {
		tok = os.Args[1]
	}

	if (tok == "") {
		log.Fatal("Error: Not enough arguments\nUsage: pytho <bot token>\n")
	}

	return tok
}

// Pytho entry point.
func main() {
	tok := token()

	_, debug := os.LookupEnv("PYTHO_DEBUG")

	var pytho Pytho
	pytho.Debug = debug

	err := pytho.Init(tok, 60)
	if err != nil {
		log.Panic(err)
	}

	pytho.Listen()
}
