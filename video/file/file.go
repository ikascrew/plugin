package file

import (
	"encoding/json"
	"fmt"
	"strings"

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

func New(param string) (*File, error) {

	p := parseParams(param)
	if p.Path == "" {
		return nil, fmt.Errorf("New File argument error")
	}

	name := p.Path
	f := File{
		name: name,
	}

	var err error

	f.cap, err = gocv.VideoCaptureFile(name)
	if err != nil {
		return nil, err
	}

	if f.cap == nil {
		return nil, fmt.Errorf("New Capture Error:[%s]", name)
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
