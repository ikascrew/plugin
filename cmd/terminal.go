package main

import (
	"fmt"

	"github.com/ikascrew/plugin/terminal"
)

func main() {
	v, err := terminal.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
