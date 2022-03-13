package main

import (
	"os"

	"github.com/rubenkristian/fakeapi/parser"
	"github.com/rubenkristian/fakeapi/server"
)

func main() {
	RunApplication()
}

func RunApplication() {
	args := os.Args

	argc := len(args)

	var file string

	if argc > 1 {
		file = args[1]
	} else {
		file = "fake.json"
	}

	p := parser.SetParser(file)

	p.RunParser()

	s := server.SetServer(p)

	s.StartServer("8080")
}
