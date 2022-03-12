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

	file := args[1]

	p := parser.SetParser(file)

	p.RunParser()

	s := server.SetServer(p)

	s.StartServer("8080")
}
