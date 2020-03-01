package server

import (
	"fmt"
	"github.com/ikascrew/core"

	cd "github.com/ikascrew/plugin/countdown"

	file "github.com/ikascrew/plugin/file"

	img "github.com/ikascrew/plugin/image"

	terminal "github.com/ikascrew/plugin/terminal"
)

var NotFoundError = fmt.Errorf("NotFound Video Type")

func Get(t string, params ...string) (core.Video, error) {

	var v video.Video
	var err error

	switch t {

	case "cd":
		v, err = cd.New(params)

	case "file":
		v, err = file.New(params)

	case "img":
		v, err = img.New(params)

	case "terminal":
		v, err = terminal.New(params)

	}

	if v == nil {
		err = NotFoundError
	}

	if err != nil {
		return nil, err
	}
	return v, nil
}
