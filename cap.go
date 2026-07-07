//go:build ignore
// TODO: 設計スケッチ(未完成)。ビルド対象外

package main

import (
	"gocv.io/x/gocv"
)

func main() {

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic(err)
	}

	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}
