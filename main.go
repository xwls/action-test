package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("hello world")
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)
	key1Bytes, err := ioutil.ReadFile("key1.txt")
	if err != nil {
		log.Panicln(err)
		return
	}
	key1 := string(key1Bytes)
	key1 = strings.TrimSpace(key1)

	key2Bytes, err := ioutil.ReadFile("key2.txt")
	if err != nil {
		log.Panicln(err)
		return
	}
	key2 := string(key2Bytes)
	key2 = strings.TrimSpace(key2)
	fmt.Println(len(key1), key1)
	fmt.Println(len(key2), key2)
}
