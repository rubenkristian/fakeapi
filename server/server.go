package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
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
			if s.fetchURLAndMatch(strings.Trim(url, "/"), strings.Trim(u, "/")) {
				fmt.Fprint(*w, result)
				return
			}
		}
		fmt.Fprint(*w, *s.parser.NotFound)
		return
	}
	fmt.Fprint(*w, "Error matching route map")
}

func (s *Server) fetchURLAndMatch(urlReq string, urlMap string) bool {
	urlReqSlices := strings.Split(urlReq, "/")
	urlMapSlices := strings.Split(urlMap, "/")

	lenURLReq := len(urlReqSlices)
	lenURLMap := len(urlMapSlices)

	if lenURLReq == lenURLMap {
		var mapPart string
		var reqPart string
		anyReg := regexp.MustCompile(`\w`)
		numberReg := regexp.MustCompile(`\d`)

		for i := 0; i < lenURLMap; i++ {
			mapPart = urlMapSlices[i]
			reqPart = urlReqSlices[i]
			if isTag(mapPart) {
				if (mapPart == "<any>" && !anyReg.MatchString(reqPart)) || (mapPart == "<number>" && !numberReg.MatchString(reqPart)) {
					return false
				}
			} else {
				if mapPart != reqPart {
					return false
				}
			}
		}

		return true
	}

	return false
}

func isTag(param string) bool {
	return param == "<any>" || param == "<number>"
}
