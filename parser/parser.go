package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type JsonFile struct {
	Services []Obj                   `json:"service"`
	NotFound *map[string]interface{} `json:"notfound"`
}

type QueryProp struct {
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Obj struct {
	Path   string                 `json:"path"`
	Method string                 `json:"method"`
	Query  map[string]QueryProp   `json:"query"`
	Result map[string]interface{} `json:"result"`
}

type Fields struct {
	Query    map[string]QueryProp
	Response map[string]interface{}
}

type Parser struct {
	service   []Obj
	MapGET    *map[string]Fields
	MapPOST   *map[string]Fields
	MapPUT    *map[string]Fields
	MapDELETE *map[string]Fields
	NotFound  *map[string]interface{}
}

func SetParser(file string) *Parser {
	jsonFile, err := os.Open("./" + file)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result JsonFile
	json.Unmarshal([]byte(byteValue), &result)

	return &Parser{
		service:  result.Services,
		NotFound: result.NotFound,
	}
}

func (p *Parser) RunParser() {
	var getList = make(map[string]Fields)
	var postList = make(map[string]Fields)
	var putList = make(map[string]Fields)
	var deleteList = make(map[string]Fields)

	for _, obj := range p.service {
		switch obj.Method {
		case "GET":
			getList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
			continue
		case "POST":
			postList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
			continue
		case "PUT":
			putList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
			continue
		case "DELETE":
			deleteList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
			continue
		}
	}

	p.MapGET = &getList
	p.MapPOST = &postList
	p.MapPUT = &putList
	p.MapDELETE = &deleteList
}

type ArrayField struct {
	Length int
	Name   string
}

func (p *Parser) isArray(field string) (bool, int, string) {
	index := 0
	lenChars := len(field)
	size := ""

	if field[index] == '[' {

		index++

		for index < lenChars {
			if field[index] == ']' {
				count, err := strconv.Atoi(size)
				if err != nil {
					log.Fatal(err)
				}
				index++
				return true, count, string(field[index:lenChars])
			} else if p.isNumber(field[index]) {
				size += string(field[index])
			}
			index++
		}
	}

	return false, 0, ""
}

func (p *Parser) isNumber(character byte) bool {
	numbers := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	for _, c := range numbers {
		if character == c {
			return true
		}
	}

	return false
}

func (p *Parser) isFunction(value string) (bool, []string) {
	arrVal := strings.Split(value, " ")
	lenVal := len(arrVal[0])

	functionList := []string{"lorem", "fullname", "firstname", "lastname"}

	if value[0] == ':' {
		funcName := value[1:(lenVal - 1)]
		for _, function := range functionList {
			if function == funcName {
				return true, arrVal[1:]
			}
		}
	}
	return false, nil
}

func (p *Parser) QueryValidation(field Fields, queries *url.Values) bool {
	for index, query := range field.Query {
		if !query.Required {
			continue
		} else if !queries.Has(index) {
			return false
		}
	}
	return true
}
