package plugin

import "testing"

// versions は未実装(TODO)の内部状態だが、init 時点の初期値が
// 崩れていないかだけ確認しておく
func TestVersionsInitialState(t *testing.T) {
	if len(versions) != 0 {
		t.Errorf("versions initial length = %d, want 0", len(versions))
	}
	if cap(versions) != 10 {
		t.Errorf("versions initial capacity = %d, want 10", cap(versions))
	}
}

// dispose は未実装(no-op)だが、呼び出してもパニックしないことを確認する
func TestDisposeNoop(t *testing.T) {
	dispose()
}
