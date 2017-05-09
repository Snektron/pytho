package main

import (
	"os"
	"fmt"
	"log"
)

func getToken() string {
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
	tok := getToken()

	var pytho Pytho
	if err := pytho.init(tok, 60); err != nil {
		log.Panic(err)
	}

	pytho.Listen()
}
