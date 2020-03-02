package terminal

import (
	"fmt"
	"image"
	"image/color"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"gocv.io/x/gocv"
)

func init() {
}

type Terminal struct {
	lines []string

	now int
	max int
}

func New(params ...string) (*Terminal, error) {

	f := Terminal{}
	var err error

	cs, err := cpu.Info()
	cpuLine := make([]string, 0)
	if err == nil {
		c := cs[0]
		cpuLine = append(cpuLine, fmt.Sprintf("    CPU -> %s x %d x %d", c.ModelName, c.Cores, len(cs)))
	} else {
		cpuLine = append(cpuLine, fmt.Sprintf("    CPU Error :%s ", err.Error()))
	}

	memLine := make([]string, 0)
	m, err := mem.VirtualMemory()
	if err == nil {
		// structが返ってきます。
		memLine = append(memLine, fmt.Sprintf("    Mem:Total: %v, Free:%v", m.Total, m.Free))
	} else {
		memLine = append(memLine, fmt.Sprintf("    Mem Error :%s ", err.Error()))
	}

	dispLine := make([]string, 0)
	dispLine = append(dispLine, fmt.Sprintf("    DISPLAY:%d x %d", 1280, 720))

	f.lines = make([]string, 8+len(cpuLine)+len(memLine)+len(dispLine))

	//CPU
	//MEM
	f.lines[0] = "I am ikascrew."
	f.lines[1] = "I am a program born to transform \"VJ System\"."
	f.lines[2] = ""
	f.lines[3] = "Today's system:"

	idx := 4
	for _, line := range cpuLine {
		f.lines[idx] = line
		idx++
	}

	for _, line := range memLine {
		f.lines[idx] = line
		idx++
	}

	for _, line := range dispLine {
		f.lines[idx] = line
		idx++
	}

	f.lines[idx] = ""

	f.lines[idx+1] = "I am a ready."
	f.lines[idx+2] = "When you're ready?"
	f.lines[idx+3] = "Let's get started!"

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
	return nil
}
