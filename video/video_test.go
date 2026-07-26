package video_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/ikascrew/plugin/video"

	"gocv.io/x/gocv"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", "file"},
		{"file", "file"},
		{"video", "file"},
		{"  file  ", "file"},
		{"FILE", "file"},
		{"img", "img"},
		{"image", "img"},
		{"IMAGE", "img"},
		{"cd", "cd"},
		{"countdown", "cd"},
		{"COUNTDOWN", "cd"},
		{"terminal", "terminal"},
		{" terminal ", "terminal"},
		{"unknown-type", "unknown-type"},
		{"  Unknown-Type  ", "unknown-type"},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got := video.Normalize(c.in)
			if got != c.want {
				t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestTypes(t *testing.T) {
	got := video.Types()
	want := []string{"file", "img", "cd", "terminal"}

	if len(got) != len(want) {
		t.Fatalf("Types() length = %d, want %d (%v)", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("Types()[%d] = %q, want %q", i, got[i], w)
		}
	}

	// 正語彙はすべて Normalize の不動点である必要がある
	for _, ty := range got {
		if video.Normalize(ty) != ty {
			t.Errorf("Normalize(%q) = %q, want fixed point", ty, video.Normalize(ty))
		}
	}
}

func TestIsGenerative(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"cd", true},
		{"countdown", true},
		{"terminal", true},
		{"file", false},
		{"video", false},
		{"img", false},
		{"image", false},
		{"", false},
		{"unknown", false},
	}

	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			if got := video.IsGenerative(c.in); got != c.want {
				t.Errorf("IsGenerative(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestSpec(t *testing.T) {
	for _, ty := range video.Types() {
		t.Run(ty, func(t *testing.T) {
			fields := video.Spec(ty)
			if fields == nil {
				t.Errorf("Spec(%q) = nil, want non-nil form definition", ty)
			}
			for _, f := range fields {
				if f.Name == "" {
					t.Errorf("Spec(%q) has field with empty Name: %+v", ty, f)
				}
			}
		})
	}

	// 旧語彙経由でも Normalize されて解決できる
	if video.Spec("countdown") == nil {
		t.Errorf("Spec(%q) = nil via legacy alias, want non-nil", "countdown")
	}

	// 未知の型は nil(ikasbox 側の JSON 直接入力フォールバック)
	if got := video.Spec("unknown"); got != nil {
		t.Errorf("Spec(unknown) = %+v, want nil", got)
	}
}

func TestGetUnknownType(t *testing.T) {
	_, err := video.Get("no-such-type", "")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
	if !errors.Is(err, video.NotFoundError) {
		t.Errorf("error should wrap video.NotFoundError, got: %v", err)
	}
}

func TestGetTerminal(t *testing.T) {
	v, err := video.Get("terminal", `{"text":"hello"}`)
	if err != nil {
		t.Fatalf("Get(terminal) error: %v", err)
	}
	defer v.Release()

	m, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if m.Empty() {
		t.Fatal("terminal frame is empty")
	}
}

func TestGetCountdownLegacyType(t *testing.T) {
	// 旧語彙 "countdown" でも解決できる
	v, err := video.Get("countdown", "")
	if err != nil {
		t.Fatalf("Get(countdown) error: %v", err)
	}
	defer v.Release()
}

func TestGetFileEmptyPathError(t *testing.T) {
	_, err := video.Get("file", `{"path":""}`)
	if err == nil {
		t.Fatal("expected error for empty file path")
	}
}

func TestGetImage(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "sample.png")

	m := gocv.NewMatWithSize(8, 8, gocv.MatTypeCV8UC3)
	defer m.Close()
	if ok := gocv.IMWrite(imgPath, m); !ok {
		t.Fatalf("failed to write test image to %s", imgPath)
	}

	v, err := video.Get("img", `{"path":"`+filepath.ToSlash(imgPath)+`"}`)
	if err != nil {
		t.Fatalf("Get(img) error: %v", err)
	}
	defer v.Release()

	frame, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if frame.Cols() != 8 || frame.Rows() != 8 {
		t.Errorf("frame size = %dx%d, want 8x8", frame.Cols(), frame.Rows())
	}

	// legacy "image" type alias も同様に解決できる
	v2, err := video.Get("image", `{"path":"`+filepath.ToSlash(imgPath)+`"}`)
	if err != nil {
		t.Fatalf("Get(image) error: %v", err)
	}
	defer v2.Release()
}

func TestGetImageMissingFileError(t *testing.T) {
	_, err := video.Get("img", `{"path":"does-not-exist.png"}`)
	if err == nil {
		t.Fatal("expected error for missing image file")
	}
}
