// Package param はプラグインが自己申告する JSON param の入力フォーム定義。
//
// ikasbox のコンテンツ登録 UI は、型名から得た []Field をそのまま JSON で
// 受け取ってフォームを組み立てる。**param の中身を解釈するのは各プラグインの
// New だけ**という原則(plugin/CLAUDE.md)を UI 側にも通すための仕組みで、
// ikasbox は器を作るだけでフィールドの意味を知らない。
//
// 型名から Field を引く入口は video.Spec。未知の型では nil が返り、
// UI 側は生の JSON 入力にフォールバックする。
package param

// Type は入力欄の種類。ikasbox の UI 側の語彙と一対一で対応するため、
// **この3つを増やさないこと**。増やすと UI 側に対応するコンポーネントが
// 無く、フォーム生成が壊れる(新しい入力が必要になった場合は、まず UI 側に
// コンポーネントを足してから追加する)
type Type string

const (
	// Text は1行のテキスト入力
	Text Type = "text"
	// Multiline は複数行のテキスト入力
	Multiline Type = "multiline"
	// DateTime は日時入力。値は RFC3339 または "2006-01-02 15:04:05"(JST)
	DateTime Type = "datetime"
)

// Field は入力欄1つ分の定義。JSON のキー名は ikasbox の UI が
// 直接参照するため変更しないこと
type Field struct {
	// Name は JSON param 内でのキー名。プラグインの Params 構造体の
	// json タグと一致させること
	Name string `json:"name"`
	// Type は入力欄の種類
	Type Type `json:"type"`
	// Label は UI に表示する見出し
	Label string `json:"label"`
	// Required は必須入力か
	Required bool `json:"required"`
	// Default は初期値。空なら JSON に出力しない
	Default string `json:"default,omitempty"`
}
