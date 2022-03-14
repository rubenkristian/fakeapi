package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type JsonFile struct {
	Services []Obj                   `json:"service"`
	NotFound *map[string]interface{} `json:"notfound"`
}

type Obj struct {
	Path   string                 `json:"path"`
	Method string                 `json:"method"`
	Query  map[string]string      `json:"query"`
	Result map[string]interface{} `json:"result"`
}

type Fields struct {
	Query    map[string]string
	Response map[string]interface{}
}

type Parser struct {
	service   []Obj
	NotFound  *map[string]interface{}
	MapGET    *map[string]Fields
	MapPOST   *map[string]Fields
	MapPUT    *map[string]Fields
	MapDELETE *map[string]Fields
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
		if obj.Method == "GET" {
			getList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
		} else if obj.Method == "POST" {
			postList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
		} else if obj.Method == "PUT" {
			putList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
		} else if obj.Method == "DELETE" {
			deleteList[obj.Path] = Fields{
				Query:    obj.Query,
				Response: obj.Result,
			}
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
