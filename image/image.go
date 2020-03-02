package image

import (
	"fmt"

	"gocv.io/x/gocv"
)

func init() {
}

type Image struct {
	name string
	src  *gocv.Mat
}

func New(params ...string) (*Image, error) {

	img := Image{
		name: params[0],
	}

	wk := gocv.IMRead(img.name, gocv.IMReadColor)
	if wk.Empty() {
		return nil, fmt.Errorf("Error:LoadImage[%s]", params[0])
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
