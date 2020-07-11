package plugin

import (
	"fmt"
	"plugin"

	"github.com/ikascrew/core"

	"golang.org/x/xerrors"
)

type videoFunc func(string, ...string) (core.Video, error)

var getVideo videoFunc

func Load() error {
	//TODO versions
	fn, err := load("plugin.so")
	if err != nil {
		return xerrors.Errorf("plugin load error: %w", err)
	}
	getVideo = fn
	return nil
}

func load(n string) (videoFunc, error) {

	p, err := plugin.Open(n)
	if err != nil {
		return nil, xerrors.Errorf("plugin open: %w", err)
	}

	sym, err := p.Lookup("GetVideo")
	if err != nil {
		return nil, xerrors.Errorf("not found GetVideo function: %w", err)
	}

	fn, ok := sym.(videoFunc)
	if !ok {
		return nil, xerrors.Errorf("GetVideo is not videoFunc: %w", err)
	}
	return fn, nil
}

func Get(f string, args ...string) (core.Video, error) {
	if getVideo == nil {
		return nil, fmt.Errorf("is not Load GetVideo function")
	}
	return getVideo(f, args...)
}
