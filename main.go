package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
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

	key2Bytes, err := ioutil.ReadFile("key2.txt")
	if err != nil {
		log.Panicln(err)
		return
	}
	key2 := string(key2Bytes)
	fmt.Println(key1 + ", " + key2)
}
