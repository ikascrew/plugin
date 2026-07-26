// Package output は生成型プラグイン(cd / terminal)の描画先解像度を保持する。
//
// 実体ファイルを持つ file / img と違い、生成型は自分で描画先を確保するため
// プロジェクトの解像度を知る必要がある。server が work file を読んだ時点と
// 同居モードでのホットリロード時に Set し、各プラグインは Size で参照する。
//
// 注意: このパッケージは 2026-07-26 の作業ツリー消失後に、呼び出し側
// (server/ikasbox.go の output.Set、countdown/terminal の output.Size)から
// API を逆算して再構成したもの。消失前の実装とは細部が異なる可能性がある。
package output

import "sync"

// 既定値は旧来のハードコード解像度(gocv.NewMatWithSize(720, 1280, ...))に合わせる
const (
	DefaultWidth  = 1280
	DefaultHeight = 720
)

var (
	mu     sync.RWMutex
	width  = DefaultWidth
	height = DefaultHeight
)

// Set は描画先解像度を設定する。0 以下の値は無視して既定値を保つ
func Set(w, h int) {
	if w <= 0 || h <= 0 {
		return
	}
	mu.Lock()
	width, height = w, h
	mu.Unlock()
}

// Size は現在の描画先解像度を返す
func Size() (int, int) {
	mu.RLock()
	defer mu.RUnlock()
	return width, height
}
