package terminal

import (
	"image"
	"image/color"
	"strings"

	"gocv.io/x/gocv"
)

func init() {
}

type Terminal struct {
	lines []string
	old   *gocv.Mat

	now int
	max int
}

func New(params ...string) (*Terminal, error) {

	f := Terminal{}

	buf := params[0]

	f.lines = strings.Split(buf, "\n")

	f.now = 0
	f.max = 0
	for _, line := range f.lines {
		f.max += len(line)
	}

	return &f, nil
}

func (v *Terminal) Next() (*gocv.Mat, error) {

	left := 20
	height := 30
	fps := 4

	//終了文字数
	n := v.now / fps

	newV := gocv.NewMatWithSize(720, 1280, gocv.MatTypeCV8UC3)

	for idx, line := range v.lines {

		buf := line

		charnum := len(line)

		if n < charnum {
			buf = line[0:n] + "|"
		}

		n -= len(line)

		gocv.PutText(&newV, buf, image.Pt(left, (idx+1)*height),
			gocv.FontHersheyComplexSmall, 1.0, color.RGBA{0, 255, 0, 0}, 2)

		//calet
		if n <= 0 {
			break
		}
	}

	v.old.Close()
	v.old = &newV

	v.now++
	return &newV, nil
}

func (v *Terminal) Wait() float64 {
	return 33.3
}

func (v *Terminal) Set(f int) {
}

func (v *Terminal) Current() int {
	return 1
}

func (v *Terminal) Source() string {
	//TODO
	return "文字列だね"
}

func (v *Terminal) Release() error {
	v.old.Close()
	return nil
}
