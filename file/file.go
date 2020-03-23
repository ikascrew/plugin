package file

import (
	"fmt"

	"gocv.io/x/gocv"
)

func init() {
}

type File struct {
	fps    float64
	frames int
	name   string
	source *gocv.Mat

	cap *gocv.VideoCapture
}

func New(args ...string) (*File, error) {

	if args == nil || len(args) == 0 {
		return nil, fmt.Errorf("New File argument error")
	}

	f := File{
		name: args[0],
	}

	var err error

	f.cap, err = gocv.VideoCaptureFile(args[0])
	if err != nil {
		return nil, err
	}

	if f.cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", f)
	}

	f.frames = int(f.cap.Get(gocv.VideoCaptureFrameCount))
	v := gocv.NewMatWithSize(720, 1280, gocv.MatTypeCV8UC3)

	f.fps = f.cap.Get(gocv.VideoCaptureFPS)

	f.source = &v
	return &f, nil
}

func (v *File) Next() (*gocv.Mat, error) {

	if v.cap == nil {
		return nil, fmt.Errorf("Error:Caputure is nil")
	}

	v.cap.Read(v.source)

	pos := int(v.cap.Get(gocv.VideoCapturePosFrames))
	if pos == v.frames {
		v.Set(1)
	}

	return v.source, nil
}

func (v *File) Wait() float64 {
	return 1000.0 / v.fps
}

func (v *File) Set(f int) {
	v.cap.Set(gocv.VideoCapturePosFrames, float64(f))
}

func (v *File) Current() int {
	return int(v.cap.Get(gocv.VideoCapturePosFrames))
}

func (v *File) Source() string {
	return v.name
}

func (v *File) Release() error {
	if v.cap != nil {
		v.cap.Close()
	}
	v.cap = nil
	return nil
}
