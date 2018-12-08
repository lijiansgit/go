package main

import (
	"flag"
	"github.com/tidwall/gjson"
	"fmt"
	"io/ioutil"
)

var (
	fileName string
	key string
)

func init() {
	flag.StringVar(&fileName, "f", "test.json", "json file")
	flag.StringVar(&key, "k", "first", "json key")
}

func main() {
	flag.Parse()
	//fmt.Println("read json file:", fileName)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	str := string(b)
	//fmt.Println("read json content:", str)

	value := gjson.Get(str, key)
	fmt.Printf("%s", value)
}


