//go:build ignore
// TODO: 設計スケッチ(未完成)。ビルド対象外

package main

import (
	"fmt"

	"github.com/ikascrew/core"

	cd "github.com/ikascrew/plugin/video/countdown"

	file "github.com/ikascrew/plugin/video/file"

	img "github.com/ikascrew/plugin/video/image"

	terminal "github.com/ikascrew/plugin/video/terminal"
)

var NotFoundError = fmt.Errorf("NotFound Video Type")

func GetPluginNames() []string {
	return []string{"cd", "file", "img", "terminal"}
}

func GetVideo(t string, params ...string) (core.Video, error) {

	var v core.Video
	var err error

	switch t {

	case "cd":
		v, err = cd.New(params...)

	case "file":
		v, err = file.New(params...)

	case "img":
		v, err = img.New(params...)

	case "terminal":
		v, err = terminal.New(params...)

	}

	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, NotFoundError
	}

	return v, nil
}

func GetEffect(t string, params ...string) (core.Effect, error) {
}

func GetTransition(t string, params ...string) (core.Effect, error) {
}
