//go:build ignore
// TODO: 設計スケッチ(未完成)。ビルド対象外

package main

import (
	"plugin"

	"github.com/ikascrew/core"
	"github.com/ikascrew/core/window"
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

	win, err := window.New("plugins test")
	if err != nil {
		panic(err)
	}

	err = win.Play(v)
	if err != nil {
		panic(err)
	}
}
