package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Service struct {
	Services []Obj `json:"service"`
}

type Obj struct {
	Path   string                 `json:"path"`
	Method string                 `json:"method"`
	Result map[string]interface{} `json:"result"`
}

type Parser struct {
	service   []Obj
	MapGET    *map[string]interface{}
	MapPOST   *map[string]interface{}
	MapPUT    *map[string]interface{}
	MapDELETE *map[string]interface{}
}

func SetParser(file string) *Parser {
	jsonFile, err := os.Open("./" + file)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result Service
	json.Unmarshal([]byte(byteValue), &result)

	return &Parser{
		service: result.Services,
	}
}

func (p *Parser) RunParser() {
	var getList = make(map[string]interface{})
	var postList = make(map[string]interface{})
	var putList = make(map[string]interface{})
	var deleteList = make(map[string]interface{})

	for _, obj := range p.service {
		if obj.Method == "GET" {
			getList[obj.Path] = obj.Result
		} else if obj.Method == "POST" {
			postList[obj.Path] = obj.Result
		} else if obj.Method == "PUT" {
			putList[obj.Path] = obj.Result
		} else if obj.Method == "DELETE" {
			deleteList[obj.Path] = obj.Result
		}
	}

	p.MapGET = &getList
	p.MapPOST = &postList
	p.MapPUT = &putList
	p.MapDELETE = &deleteList
}
