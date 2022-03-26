package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/rubenkristian/fakeapi/parser"
)

type Server struct {
	parser *parser.Parser
}

func SetServer(parser *parser.Parser) *Server {
	return &Server{
		parser: parser,
	}
}

func (s *Server) StartServer(port string) {
	var httpServerError = make(chan error)
	var waitGroup sync.WaitGroup

	var address string = ":" + port

	listen, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", s.serviceHandler)
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		httpServerError <- http.Serve(listen, nil)
	}()

	select {
	case err := <-httpServerError:
		fmt.Println("Service could not be started", err)
	default:
		fmt.Println("FakeAPI serve on http://localhost:" + port)
	}

	waitGroup.Wait()
}

func (s *Server) serviceHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	method := req.Method
	query := req.URL.Query()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.matchURL(&w, &query, url, method)
}
