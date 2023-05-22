package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rubenkristian/fakeapi/parser"
	"github.com/valyala/fasthttp"
)

func (s *Server) fastMatchURL(ctx *fasthttp.RequestCtx, url string, method string) {
	p := s.parser

	var mapMethod *map[string]parser.Fields

	switch method {
	case "GET":
		mapMethod = p.MapGET
	case "POST":
		mapMethod = p.MapPOST
	case "PUT":
		mapMethod = p.MapPUT
	case "DELETE":
		mapMethod = p.MapDELETE
	}

	if mapMethod != nil {
		for u, result := range *mapMethod {
			isMatch, params := s.fetchURLAndMatch(strings.Trim(url, "/"), strings.Trim(u, "/"))
			// TODO: fetch result and create anonym function to give result and response to requester
			if isMatch {
				// TODO: do something with params
				var status = result.Response["status"].(float64)
				var response = result.Response["response"]
				var byteJson, err = json.Marshal(response)

				if err != nil {
					ctx.Response.SetBodyString(err.Error())
				} else {
					var stringJson = string(byteJson)
					(*ctx).Response.Header.Set("Content-Type", "application/json; charset=utf-8")
					if status == 200 {
						(*ctx).Response.Header.SetStatusCode(http.StatusOK)
					}

					for index, param := range params {
						var indexParam = fmt.Sprintf("${%d}", index)
						fmt.Print(indexParam)
						stringJson = strings.ReplaceAll(stringJson, indexParam, param)
					}

					(*ctx).Response.SetBodyString(stringJson)
				}
				return
			}
		}

		(*ctx).Response.SetStatusCode(http.StatusNotFound)
		(*ctx).Response.SetBodyString("Route not found")
		return
	}

	(*ctx).Response.SetStatusCode(http.StatusNotFound)
	(*ctx).Response.SetBodyString("Route not found")
}

func (s *Server) matchURL(w *http.ResponseWriter, query *url.Values, url string, method string) {
	p := s.parser

	var mapMethod *map[string]parser.Fields

	switch method {
	case "GET":
		mapMethod = p.MapGET
	case "POST":
		mapMethod = p.MapPOST
	case "PUT":
		mapMethod = p.MapPUT
	case "DELETE":
		mapMethod = p.MapDELETE
	}

	if mapMethod != nil {
		for u, result := range *mapMethod {
			isMatch, params := s.fetchURLAndMatch(strings.Trim(url, "/"), strings.Trim(u, "/"))
			// TODO: fetch result and create anonym function to give result and response to requester
			if isMatch {
				// TODO: do something with params
				var status = result.Response["status"].(float64)
				var response = result.Response["response"]
				var byteJson, err = json.Marshal(response)

				if err != nil {
					fmt.Fprint(*w, err.Error())
				} else {
					var stringJson = string(byteJson)
					(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
					if status == 200 {
						(*w).WriteHeader(http.StatusOK)
					}

					for index, param := range params {
						var indexParam = fmt.Sprintf("${%d}", index)
						fmt.Print(indexParam)
						stringJson = strings.ReplaceAll(stringJson, indexParam, param)
					}

					fmt.Fprint(*w, string(stringJson))
				}
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
