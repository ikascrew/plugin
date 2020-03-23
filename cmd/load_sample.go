package main

import (
	"fmt"
	"plugin"

	"github.com/ikascrew/core"
)

func main() {

	p, err := plugin.Open("plugins.so")
	if err != nil {
		panic(err)
	}

	sym, err := p.Lookup("GetVideo")
	if err != nil {
		panic(err)
	}

	fn, ok := sym.(func(string, ...string) (core.Video, error))
	if !ok {
		panic("not cast GetVideo function")
	}

	v, err := fn("file", "newrising.mp4")
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
