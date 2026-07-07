// Package video は ikascrew の映像プラグインのレジストリ。
//
// コンテンツの type 文字列(正語彙: file / img / cd / terminal)から
// core.Video 実装を生成する唯一の対応表。ikasbox / server は
// この Get を通して型を解決し、型ごとの分岐を自前で持たない。
//
// param は JSON 文字列で、解釈は各プラグインの New だけが行う。
// 中間層(ikasbox の DB・server の work file・gRPC)は不透明な
// 文字列として運ぶこと。JSON でない場合は各プラグインが
// 旧来の生文字列(パスやテキスト)として解釈する。
package video

import (
	"strings"

	"github.com/ikascrew/core"

	cd "github.com/ikascrew/plugin/video/countdown"
	file "github.com/ikascrew/plugin/video/file"
	img "github.com/ikascrew/plugin/video/image"
	terminal "github.com/ikascrew/plugin/video/terminal"

	"golang.org/x/xerrors"
)

var NotFoundError = xerrors.New("NotFound Video Type")

// Types は正語彙の一覧を返す
func Types() []string {
	return []string{"file", "img", "cd", "terminal"}
}

// IsGenerative は実体ファイルを持たない生成型プラグインかを返す。
// ikasbox でのコンテンツ登録方法(import か params 登録か)の判定に使う
func IsGenerative(t string) bool {
	switch Normalize(t) {
	case "cd", "terminal":
		return true
	}
	return false
}

// Normalize は各所で揺れている型名を正語彙に寄せる。
// ikasbox の旧データ("image")や client の旧指定("countdown")を吸収する
func Normalize(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "", "file", "video":
		return "file"
	case "img", "image":
		return "img"
	case "cd", "countdown":
		return "cd"
	case "terminal":
		return "terminal"
	}
	return strings.ToLower(strings.TrimSpace(t))
}

// Get は型名と JSON param から Video を生成する
func Get(t string, param string) (core.Video, error) {

	var v core.Video
	var err error

	switch Normalize(t) {
	case "file":
		v, err = file.New(param)
	case "img":
		v, err = img.New(param)
	case "cd":
		v, err = cd.New(param)
	case "terminal":
		v, err = terminal.New(param)
	default:
		return nil, xerrors.Errorf("video type[%s]: %w", t, NotFoundError)
	}

	if err != nil {
		return nil, xerrors.Errorf("video new[%s]: %w", t, err)
	}
	return v, nil
}
