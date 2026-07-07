//go:build ignore
// TODO: 設計スケッチ(未完成)。ビルド対象外

package main

import (
	"github.com/ikascrew/core/window"
	"github.com/ikascrew/plugin/video/telop"
)

func main() {

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
