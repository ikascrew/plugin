package plugin

import "testing"

// Get は Load が成功していない(getVideo が未設定)状態ではエラーを返す。
// このリポジトリの動的ロード(Go plugin)は Linux 限定の設計スケッチであり、
// このテスト環境(Windows)では Load を成功させられないため、
// 未ロード時のエラーハンドリングだけを検証する
func TestGetWithoutLoad(t *testing.T) {
	if getVideo != nil {
		t.Skip("getVideo already set by another test; skipping to avoid interference")
	}

	_, err := Get("file", "sample.mp4")
	if err == nil {
		t.Fatal("expected error when GetVideo function has not been loaded")
	}
}

// 存在しない plugin.so の Load はエラーを返し、getVideo を書き換えない
func TestLoadMissingFile(t *testing.T) {
	before := getVideo

	if err := Load(); err == nil {
		t.Fatal("expected error loading a nonexistent plugin.so")
	}

	if getVideo != nil && before == nil {
		t.Error("Load should not set getVideo on failure")
	}
}
