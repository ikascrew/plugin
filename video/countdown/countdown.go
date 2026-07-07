package countdown

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"strings"
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

// Params は JSON param の形。解釈はこのプラグインだけが行う。
// Target は RFC3339("2026-12-31T23:59:59+09:00")または
// "2006-01-02 15:04:05"(JST 扱い)。Text はカウント終了後に表示する文字列
type Params struct {
	Target string `json:"target"`
	Text   string `json:"text"`
}

func parseParams(param string) Params {
	p := Params{}
	s := strings.TrimSpace(param)
	if strings.HasPrefix(s, "{") {
		if err := json.Unmarshal([]byte(s), &p); err == nil {
			return p
		}
	}
	// JSON でなければ旧来どおり終了後テキストとして扱う
	p.Text = param
	return p
}

func parseTarget(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, jst)
}

func New(param string) (*Countdown, error) {

	p := parseParams(param)

	f := Countdown{}
	f.text = p.Text

	// target 未指定は旧来のハードコード値(過去日時 = 即終了表示)
	f.target = Target.In(jst).Unix()
	if p.Target != "" {
		t, err := parseTarget(p.Target)
		if err != nil {
			return nil, fmt.Errorf("countdown target[%s]: %v", p.Target, err)
		}
		f.target = t.Unix()
	}

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

	} else if v.text != "" {
		gocv.PutText(&mat, v.text, image.Pt(60, 400),
			gocv.FontHersheyComplexSmall, 7.4, color.RGBA{255, 255, 255, 0}, 4)
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
