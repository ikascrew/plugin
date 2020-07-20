package telop_test

import (
	"github.com/ikascrew/core/window"
	"github.com/ikascrew/plugin/telop"
)

func ExampleNew() {

	v, err := telop.New("logo.svg")
	if err != nil {
		panic(err)
	}

	win, err := window.New("file play example")
	if err != nil {
		panic(err)
	}

	err = win.Play(v)
	if err != nil {
		panic(err)
	}
}
