package image

import (
	"encoding/json"
	"fmt"
	"strings"

	"gocv.io/x/gocv"
)

func init() {
}

type Image struct {
	name string
	src  *gocv.Mat
}

// Params は JSON param の形。解釈はこのプラグインだけが行う
type Params struct {
	Path string `json:"path"`
}

func parseParams(param string) Params {
	p := Params{}
	s := strings.TrimSpace(param)
	if strings.HasPrefix(s, "{") {
		if err := json.Unmarshal([]byte(s), &p); err == nil {
			return p
		}
	}
	// JSON でなければ旧来どおりパス文字列として扱う
	p.Path = s
	return p
}

func New(param string) (*Image, error) {

	p := parseParams(param)
	if p.Path == "" {
		return nil, fmt.Errorf("New Image argument error")
	}

	img := Image{
		name: p.Path,
	}

	wk := gocv.IMRead(img.name, gocv.IMReadColor)
	if wk.Empty() {
		return nil, fmt.Errorf("Error:LoadImage[%s]", img.name)
	}

	img.src = &wk

	return &img, nil
}

func (v *Image) Next() (*gocv.Mat, error) {
	return v.src, nil
}

func (v *Image) Wait() float64 {
	return 33.3
}

func (v *Image) Set(f int) {
}

func (v *Image) Current() int {
	return 1
}

func (v *Image) Source() string {
	return v.name
}

func (v *Image) Release() error {
	if !v.src.Empty() {
		v.src.Close()
	}
	return nil
}
