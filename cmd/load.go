package main

import (
	"fmt"

	"github.com/ikascrew/plugin"
)

func main() {
	err := plugin.Load()
	if err != nil {
		panic(err)
	}

	v, err := plugin.Get("file", "inspect.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
