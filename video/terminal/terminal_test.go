package terminal_test

import (
	"testing"

	"github.com/ikascrew/plugin/video/output"
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

// New に渡す param の解釈(JSON 正常系・旧来の生文字列・不正 JSON)を検証する。
// terminal は非 JSON・不正 JSON のどちらも「全体を表示テキストとして扱う」
// フォールバックがあるため、常に New はエラーにならない
func TestNewParamVariants(t *testing.T) {
	cases := []struct {
		name  string
		param string
	}{
		{"JSON single line", `{"text":"hello"}`},
		{"JSON multi line", `{"text":"line1\nline2\nline3"}`},
		{"legacy raw string", "plain text\nmore lines"},
		{"malformed JSON falls back to raw text", `{"text":`},
		{"empty param", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v, err := terminal.New(c.param)
			if err != nil {
				t.Fatalf("New(%q) error: %v", c.param, err)
			}
			defer v.Release()

			m, err := v.Next()
			if err != nil {
				t.Fatalf("Next error: %v", err)
			}
			if m.Empty() {
				t.Errorf("frame is empty for param %q", c.param)
			}
		})
	}
}

// Next は毎フレーム新しい Mat を確保する(内部の old を都度 Close する)。
// 複数回呼んでもパニックしないこと、フレームごとに別の Mat が
// 返ってくることを確認する
func TestNextAllocatesNewMatEachCall(t *testing.T) {
	v, err := terminal.New(`{"text":"progressive reveal test"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	var prev *interface{}
	_ = prev

	m1, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	m2, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}

	if m1 == m2 {
		t.Errorf("expected distinct Mat pointers across frames, got same pointer %p", m1)
	}

	// さらに数フレーム進めてもパニック/エラーがないこと
	for i := 0; i < 10; i++ {
		if _, err := v.Next(); err != nil {
			t.Fatalf("Next error at frame %d: %v", i, err)
		}
	}
}

// Next が一度も呼ばれていない状態での Release はパニックしない
func TestReleaseWithoutNext(t *testing.T) {
	v, err := terminal.New(`{"text":"hi"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	if err := v.Release(); err != nil {
		t.Fatalf("Release error: %v", err)
	}
}

// Next はプロジェクト解像度(output.Size())のキャンバスに描画する
func TestNextRespectsOutputResolution(t *testing.T) {
	w, h := 1920, 1080
	output.Set(w, h)
	defer output.Set(output.DefaultWidth, output.DefaultHeight)

	v, err := terminal.New(`{"text":"resolution test"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	m, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if m.Cols() != w || m.Rows() != h {
		t.Errorf("frame size = %dx%d, want %dx%d", m.Cols(), m.Rows(), w, h)
	}
}

func TestSpec(t *testing.T) {
	fields := terminal.Spec()
	if len(fields) == 0 {
		t.Fatal("Spec() returned no fields")
	}
	found := false
	for _, f := range fields {
		if f.Name == "text" {
			found = true
		}
	}
	if !found {
		t.Errorf("Spec() missing 'text' field: %+v", fields)
	}
}
