package countdown

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"
)

func init() {
}

var loc, _ = time.LoadLocation("Asia/Tokyo")
var Target = time.Date(2022, time.January, 1, 0, 0, 0, 0, loc)
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type Countdown struct {
	text   string
	target int64
}

func New(params ...string) (*Countdown, error) {

	f := Countdown{}

	//最終文字列
	f.text = params[0]

	//TODO ターゲットも設定
	f.target = Target.In(jst).Unix()
	return &f, nil
}

func (v *Countdown) Next() (*gocv.Mat, error) {

	mat := gocv.NewMatWithSize(720, 1280, gocv.MatTypeCV8UC3)

	now := time.Now().In(jst)
	d := v.target - now.Unix()

	if d >= 0 {

		buf := fmt.Sprintf("%d", d)

		// 3 200
		// 2 300
		// 1 400

		left := 500 - (len(buf) * 100)

		gocv.PutText(&mat, buf, image.Pt(left, 400),
			gocv.FontHersheyComplexSmall, 16.0, color.RGBA{255, 255, 255, 0}, 4)

		if len(buf) <= 1 {
			gocv.Circle(&mat, image.Pt(502, 295), 200, color.RGBA{255, 255, 255, 0}, 8)
		}

	} else {
		gocv.PutText(&mat, "Happy", image.Pt(180, 200),
			gocv.FontHersheyComplexSmall, 9.0, color.RGBA{255, 255, 255, 0}, 4)
		gocv.PutText(&mat, "New Year!", image.Pt(10, 450),
			gocv.FontHersheyComplexSmall, 7.4, color.RGBA{255, 255, 255, 0}, 4)
	}

	return &mat, nil
}

func (v *Countdown) Wait() float64 {
	return 33.3
}

func (v *Countdown) Set(f int) {
}

func (v *Countdown) Current() int {
	return 1
}

func (v *Countdown) Source() string {
	//TODO
	return "文字列返すか？"
}

func (v *Countdown) Release() error {
	return nil
}
