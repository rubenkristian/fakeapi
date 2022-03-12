package server

import (
	"fmt"
	"net/http"

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
	var address string = ":" + port

	http.HandleFunc("/", s.serviceHandler)
	http.ListenAndServe(address, nil)
}

func (s *Server) serviceHandler(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	method := req.Method

	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.matchURL(&w, url, method)
}

func (s *Server) matchURL(w *http.ResponseWriter, url string, method string) {
	p := s.parser

	var mapMethod *map[string]interface{}

	if method == "GET" {
		mapMethod = p.MapGET
	} else if method == "POST" {
		mapMethod = p.MapPOST
	} else if method == "PUT" {
		mapMethod = p.MapPUT
	} else if method == "DELETE" {
		mapMethod = p.MapDELETE
	}

	if mapMethod != nil {
		for u, result := range *mapMethod {
			if u == url {
				fmt.Fprint(*w, result)
				break
			}
		}
	}
}
