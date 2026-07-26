package image_test

import (
	"path/filepath"
	"testing"

	"github.com/ikascrew/plugin/video/image"

	"gocv.io/x/gocv"
)

// writeTempImage は t.TempDir() 配下に gocv.IMWrite で一時画像を生成する
func writeTempImage(t *testing.T, w, h int) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.png")

	m := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC3)
	defer m.Close()

	if ok := gocv.IMWrite(path, m); !ok {
		t.Fatalf("failed to write test image to %s", path)
	}
	return path
}

func TestNewJSONPath(t *testing.T) {
	path := writeTempImage(t, 16, 12)

	v, err := image.New(`{"path":"` + filepath.ToSlash(path) + `"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	m, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if m.Cols() != 16 || m.Rows() != 12 {
		t.Errorf("frame size = %dx%d, want 16x12", m.Cols(), m.Rows())
	}
	if m.Mat().Channels() != 3 {
		t.Errorf("channels = %d, want 3", m.Mat().Channels())
	}
	if want := filepath.ToSlash(path); v.Source() != want {
		t.Errorf("Source() = %q, want %q", v.Source(), want)
	}
}

func TestNewLegacyRawPath(t *testing.T) {
	path := writeTempImage(t, 10, 10)

	// JSON でない場合は全体をパスとして扱う(旧来の生文字列)
	v, err := image.New(path)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	if v.Source() != path {
		t.Errorf("Source() = %q, want %q", v.Source(), path)
	}
}

func TestNewEmptyPathError(t *testing.T) {
	if _, err := image.New(`{"path":""}`); err == nil {
		t.Fatal("expected error for empty path")
	}
	if _, err := image.New(""); err == nil {
		t.Fatal("expected error for empty raw string")
	}
}

func TestNewMissingFileError(t *testing.T) {
	if _, err := image.New(`{"path":"does-not-exist.png"}`); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewMalformedJSONFallsBackToRawPath(t *testing.T) {
	path := writeTempImage(t, 10, 10)

	// "{" で始まるが不正な JSON。パース失敗時は全体を生パスとして扱う契約
	malformed := `{"path":"` + filepath.ToSlash(path)

	if _, err := image.New(malformed); err == nil {
		t.Fatal("expected error, malformed json string is not a valid path")
	}
}

// Next は毎回同じソース Mat を返す(フレームの再読込をしない)ことを確認する
func TestNextReturnsSameSource(t *testing.T) {
	path := writeTempImage(t, 10, 10)

	v, err := image.New(`{"path":"` + filepath.ToSlash(path) + `"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	m1, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	m2, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if m1 != m2 {
		t.Errorf("Next() should return the same Mat pointer, got %p != %p", m1, m2)
	}
}

func TestWaitSetCurrent(t *testing.T) {
	path := writeTempImage(t, 4, 4)

	v, err := image.New(`{"path":"` + filepath.ToSlash(path) + `"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	if v.Wait() <= 0 {
		t.Errorf("Wait() = %v, want > 0", v.Wait())
	}
	// Set は no-op、Current は常に 1(静止画なので)
	v.Set(10)
	if v.Current() != 1 {
		t.Errorf("Current() = %d, want 1", v.Current())
	}
}

func TestSpec(t *testing.T) {
	fields := image.Spec()
	if len(fields) == 0 {
		t.Fatal("Spec() returned no fields")
	}
	found := false
	for _, f := range fields {
		if f.Name == "path" {
			found = true
			if !f.Required {
				t.Errorf("path field should be required")
			}
		}
	}
	if !found {
		t.Errorf("Spec() missing 'path' field: %+v", fields)
	}
}
