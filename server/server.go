package server

import (
	"fmt"
	"sync"

	"github.com/rubenkristian/fakeapi/parser"
	"github.com/valyala/fasthttp"
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

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		httpServerError <- fasthttp.ListenAndServe(address, s.serviceHandler)
	}()

	select {
	case err := <-httpServerError:
		fmt.Println("Service could not be started", err)
	default:
		fmt.Println("FakeAPI serve on http://localhost:" + port)
	}

	waitGroup.Wait()
}

func (s *Server) serviceHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Request.Header.Method())

	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	s.fastMatchURL(ctx, path, method)
}

// func (s *Server) serviceHandler(w http.ResponseWriter, req *http.Request) {
// 	url := req.URL.String()
// 	method := req.Method
// 	query := req.URL.Query()

// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	s.matchURL(&w, &query, url, method)
// }
