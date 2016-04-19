package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Router struct {
	Pattern string `json: pattern`
	Repl    string `json: repl`
	Comment string `json: comment`
}

type Table struct {
	Hosts   []string `json:"hosts"`
	Routers []Router `json:"routers"`
	Version string   `json:"version"`
}

// read a urlrewrite json file
func readJson(fp string) *Table {
	file, err := ioutil.ReadFile(fp)
	var inputJson = new(Table)
	err = json.Unmarshal([]byte(string(file)), &inputJson)
	if err != nil {
		panic(err)
	}
	return inputJson
}

// transfer named group
func transferRules(inputJson *Table) []byte {
	for i, val := range inputJson.Routers {
		pattern := val.Pattern
		repl := val.Repl

		re := regexp.MustCompile(pattern)
		for i := 1; i <= re.NumSubexp(); i++ {
			repl = strings.Replace(repl, re.SubexpNames()[i], strconv.Itoa(i), -1)
		}
		inputJson.Routers[i].Repl = repl
	}

	outputJson, err := json.MarshalIndent(inputJson, "", "\t")
	if err != nil {
		panic(err)
	}

	outputJson = bytes.Replace(outputJson, []byte("\\u003c"), []byte("<"), -1)
	outputJson = bytes.Replace(outputJson, []byte("\\u003e"), []byte(">"), -1)
	outputJson = bytes.Replace(outputJson, []byte("\\u0026"), []byte("&"), -1)
	outputJson = bytes.Replace(outputJson, []byte("?P<key_prefix>"), []byte(""), -1)
	outputJson = bytes.Replace(outputJson, []byte("?P<xxx>"), []byte(""), -1)
	outputJson = bytes.Replace(outputJson, []byte("?P<yyy>"), []byte(""), -1)
	outputJson = bytes.Replace(outputJson, []byte("?P<from_format>"), []byte(""), -1)
	outputJson = bytes.Replace(outputJson, []byte("?P<to_format>"), []byte(""), -1)

	return outputJson
}

func main() {

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: transfer json_path")
		return
	}

	fp := flag.Arg(0)

	fmt.Println(string(transferRules(readJson(fp))))
}
