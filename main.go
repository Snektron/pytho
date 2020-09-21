package main

import (
	"os"
	"log"
	"io/ioutil"
	"strings"
)

// Pytho entry point.
func main() {
	// tok := token()
	if len(os.Args) < 2 {
		log.Fatal("Error: Not enough arguments\nUsage: pytho <bot token file>\n")
	}
	tok, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("Error: Couldn't read token file\n")
	}

	_, debug := os.LookupEnv("PYTHO_DEBUG")

	var pytho Pytho
	pytho.Debug = debug

	err = pytho.Init(strings.Trim(string(tok), "\n"), 60)
	if err != nil {
		log.Panic(err)
	}

	pytho.Listen()
}
