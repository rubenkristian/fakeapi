package server

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rubenkristian/fakeapi/parser"
)

func (s *Server) matchURL(w *http.ResponseWriter, query *url.Values, url string, method string) {
	p := s.parser

	var mapMethod *map[string]parser.Fields

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
			isMatch, params := s.fetchURLAndMatch(strings.Trim(url, "/"), strings.Trim(u, "/"))
			// TODO: fetch result and create anonym function to give result and response to requester
			if isMatch {
				// TODO: do something with params
				fmt.Fprint(*w, result)
				fmt.Fprint(*w, params)
				return
			}
		}

		fmt.Fprint(*w, *s.parser.NotFound)
		return
	}

	fmt.Fprint(*w, "Error matching route map")
}

func (s *Server) fetchURLAndMatch(urlReq string, urlMap string) (bool, []string) {
	urlReqSlices := strings.Split(urlReq, "/")
	urlMapSlices := strings.Split(urlMap, "/")

	lenURLReq := len(urlReqSlices)
	lenURLMap := len(urlMapSlices)

	var arrURLParam []string

	// status := true

	if lenURLReq == lenURLMap {
		var mapPart string
		var reqPart string

		anyReg := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
		numberReg := regexp.MustCompile(`^[0-9]+$`)

		for i := 0; i < lenURLMap; i++ {
			mapPart = urlMapSlices[i]
			reqPart = urlReqSlices[i]

			if isTag(mapPart) {
				if (mapPart == "<any>" && !anyReg.MatchString(reqPart)) || (mapPart == "<number>" && !numberReg.MatchString(reqPart)) {
					return false, nil
				}

				arrURLParam = append(arrURLParam, reqPart)
			} else {
				if mapPart != reqPart {
					return false, nil
				}
			}
		}
		return true, arrURLParam
	}

	return false, nil
}

func isTag(param string) bool {
	return param == "<any>" || param == "<number>"
}
