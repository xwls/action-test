package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	envKey1, b := os.LookupEnv("KEY1")
	fmt.Println(len(envKey1), envKey1, b)
	envKey2, b := os.LookupEnv("KEY2")
	fmt.Println(len(envKey2), envKey2, b)
}
