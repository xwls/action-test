package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	fmt.Println("hello world")
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)

	envKey1, b := os.LookupEnv("KEY1")
	fmt.Println(len(envKey1), envKey1, b)
	envKey2, b := os.LookupEnv("KEY2")
	fmt.Println(len(envKey2), envKey2, b)
	token, b := os.LookupEnv("GITHUB_TOKEN")
	fmt.Println(len(token), token, b)
}
