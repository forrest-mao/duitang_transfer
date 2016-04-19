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

// json format
// {
// 	"hosts": [
// 		"itisatest.qiniudn.com"
// 	],
// 	"routers": [
// 		{
// 			"Pattern": "^((?:/[^/]+)*/[^.]+)[.]thumb[.]([1-9]\\d*)_0_c[.]([^_]+)_(jpeg|webp)$",
// 			"Repl": "${1}.${3}?imageMogr2/format/${4}/quality/90/thumbnail/${2}x/ignore-error/1",
// 			"Comment": "x_0 / convert_format"
// 		},
// 		{
// 			"Pattern": "^((?:/[^/]+)*/[^.]+)[.]thumb[.]0_([1-9]\\d*)_c[.]([^_]+)_(jpeg|webp)$",
// 			"Repl": "${1}.${3}?imageMogr2/format/${4}/quality/90/thumbnail/x${2}/ignore-error/1",
// 			"Comment": "0_x / convert_format"
// 		}
// 	]
// 	"version": "71848b9ef9074fbf9c5cfec206f8e27b"
// }

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

// transfer named group routers to numbered group routers
func transferRules(inputJson *Table) []byte {

	for i, val := range inputJson.Routers {
		pattern := val.Pattern
		repl := val.Repl

		// replace named group to numbered group
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

	// replace pattern strings
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
	if flag.NArg() < 2 {
		fmt.Println("Usage: ./transfer json_from_path json_to_path")
		return
	}

	fromPath := flag.Arg(0)
	toPath := flag.Arg(1)

	//fmt.Println(string(transferRules(readJson(fp))))

	// write outputJson to file
	err := ioutil.WriteFile(toPath, transferRules(readJson(fromPath)), 0666)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("transfer to " + toPath + " finished")
	}
}
