package main

import (
	"os"
	"fmt"
	"log"
)

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

func main() {
	tok := token()

	var bot Bot
	if err := bot.Init(tok, 60); err != nil {
		log.Panic(err)
	}

	bot.Listen()
}
