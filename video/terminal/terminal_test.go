package terminal_test

import (
	"testing"

	"github.com/ikascrew/plugin/video/terminal"
)

func TestTerminal(t *testing.T) {
	v, err := terminal.New(`{"text":"hello\nterminal"}`)
	if err != nil {
		t.Fatal(err)
	}
	m, err := v.Next()
	if err != nil {
		t.Fatal(err)
	}
	if m.Rows() == 0 || m.Cols() == 0 {
		t.Fatalf("empty frame: %dx%d", m.Rows(), m.Cols())
	}
	if err := v.Release(); err != nil {
		t.Fatal(err)
	}
}
