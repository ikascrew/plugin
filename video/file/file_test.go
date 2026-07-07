package file_test

import (
	"github.com/ikascrew/core/window"
	"github.com/ikascrew/plugin/video/file"
)

func ExampleNew() {

	v, err := file.New(`{"path":"sample.mp4"}`)
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
