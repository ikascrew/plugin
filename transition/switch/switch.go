//go:build ignore
// TODO: 設計スケッチ(未完成)。ビルド対象外

package transition

import (
	"strconv"

	"gocv.io/x/gocv"
	"golang.org/x/xerrors"
)

func New(params ...string) (*Switch, error) {
	var s Switch
	var err error
	if len(params) >= 1 {
		s.max, err = strconv.ParseFloat(params[0])
	} else {
		s.max = 200.0
	}
	if err != nil {
		return nil, xerrors.Errorf("switch Effect error: %w", err)
	}
	return &s, nil
}

type Switch struct {
	max float64
	now float64
}

func (s *Switch) Set(v interface{}) error {
	s.now, ok = v.(float64)
	if !ok {
		return xerrors.Errorf("Now Value Parse(float64) error")
	}
	return nil
}

func (s *Switch) Convert(now Video, next Video, out *gocv.Mat) error {
	alpha := s.now / s.max
	nowM := now.Next()
	nextM := next.Next()
	gocv.AddWeighted(nextM, float64(alpha), *out, float64(1.0-alpha), 0.0, nowM)
	return nil
}
