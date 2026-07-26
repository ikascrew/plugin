package file_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikascrew/plugin/video/file"
)

// copySample は sample.mp4(動画代替の実ファイル)を t.TempDir() 配下へ
// コピーし、そのパスを返す
func copySample(t *testing.T) string {
	t.Helper()

	src, err := os.Open("sample.mp4")
	if err != nil {
		t.Fatalf("open sample.mp4: %v", err)
	}
	defer src.Close()

	dir := t.TempDir()
	dstPath := filepath.Join(dir, "video.mp4")

	dst, err := os.Create(dstPath)
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		t.Fatalf("copy sample: %v", err)
	}
	return dstPath
}

func TestNewJSONPath(t *testing.T) {
	v, err := file.New(`{"path":"sample.mp4"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	if v.Source() != "sample.mp4" {
		t.Errorf("Source() = %q, want %q", v.Source(), "sample.mp4")
	}
	if v.Wait() <= 0 {
		t.Errorf("Wait() = %v, want > 0", v.Wait())
	}

	m, err := v.Next()
	if err != nil {
		t.Fatalf("Next error: %v", err)
	}
	if m.Empty() {
		t.Fatal("frame is empty")
	}
}

func TestNewLegacyRawPath(t *testing.T) {
	path := copySample(t)

	// JSON でない場合は全体をパスとして扱う(旧来の生文字列)
	v, err := file.New(path)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	if v.Source() != path {
		t.Errorf("Source() = %q, want %q", v.Source(), path)
	}
}

func TestNewEmptyPathError(t *testing.T) {
	if _, err := file.New(`{"path":""}`); err == nil {
		t.Fatal("expected error for empty path")
	}
	if _, err := file.New(""); err == nil {
		t.Fatal("expected error for empty raw string")
	}
}

func TestNewMissingFileError(t *testing.T) {
	if _, err := file.New(`{"path":"does-not-exist.mp4"}`); err == nil {
		t.Fatal("expected error for missing video file")
	}
}

// フレーム末尾に達すると先頭(フレーム1)へループする挙動を検証する。
// sample.mp4 は 609 フレーム(24fps)なので、末尾付近にシークしてから
// 数フレーム進めればループを確認できる
func TestNextLoopsToStartAtEnd(t *testing.T) {
	path := copySample(t)

	v, err := file.New(`{"path":"` + filepath.ToSlash(path) + `"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}
	defer v.Release()

	v.Set(607) // 末尾の2フレーム手前へシーク

	sawWrap := false
	prev := v.Current()
	for i := 0; i < 5; i++ {
		if _, err := v.Next(); err != nil {
			t.Fatalf("Next error: %v", err)
		}
		cur := v.Current()
		if cur < prev {
			sawWrap = true
		}
		prev = cur
	}

	if !sawWrap {
		t.Errorf("expected frame position to wrap back to start, last pos=%d", prev)
	}
}

func TestReleaseIsIdempotent(t *testing.T) {
	v, err := file.New(`{"path":"sample.mp4"}`)
	if err != nil {
		t.Fatalf("New error: %v", err)
	}

	if err := v.Release(); err != nil {
		t.Fatalf("Release error: %v", err)
	}
	// 2回目の Release はパニックしない
	if err := v.Release(); err != nil {
		t.Fatalf("second Release error: %v", err)
	}

	// Release 後の Next はエラーを返す(cap が nil)
	if _, err := v.Next(); err == nil {
		t.Error("expected error calling Next after Release")
	}
}

func TestSpec(t *testing.T) {
	fields := file.Spec()
	if len(fields) == 0 {
		t.Fatal("Spec() returned no fields")
	}
	found := false
	for _, f := range fields {
		if f.Name == "path" {
			found = true
		}
	}
	if !found {
		t.Errorf("Spec() missing 'path' field: %+v", fields)
	}
}
