package param_test

import (
	"encoding/json"
	"testing"

	"github.com/ikascrew/plugin/video/param"
)

// Type 定数が UI 側の語彙("text"/"multiline"/"datetime")と一致することを
// 検証する。CLAUDE.md はこの3語彙を増やさない方針を明記しているため、
// 値そのものが変わっていないかを固定する
func TestTypeConstants(t *testing.T) {
	cases := []struct {
		name string
		typ  param.Type
		want string
	}{
		{"text", param.Text, "text"},
		{"multiline", param.Multiline, "multiline"},
		{"datetime", param.DateTime, "datetime"},
	}
	for _, c := range cases {
		if string(c.typ) != c.want {
			t.Errorf("%s: got %q, want %q", c.name, c.typ, c.want)
		}
	}
}

// Field は ikasbox の登録 UI が消費する JSON 形。キー名が変わると
// UI 側のフォーム生成が壊れるため固定する
func TestFieldJSONKeys(t *testing.T) {
	f := param.Field{
		Name:     "target",
		Type:     param.DateTime,
		Label:    "Target",
		Required: true,
		Default:  "",
	}

	buf, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	for _, key := range []string{"name", "type", "label", "required"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in %s", key, buf)
		}
	}

	// Default は omitempty なので空文字なら出力に含まれない
	if _, ok := m["default"]; ok {
		t.Errorf("default should be omitted when empty: %s", buf)
	}
}

// Default が非空の場合は出力される(omitempty の反対側の確認)
func TestFieldJSONDefaultPresent(t *testing.T) {
	f := param.Field{Name: "text", Type: param.Text, Label: "Text", Default: "hello"}

	buf, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(buf, &m); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if v, ok := m["default"]; !ok || v != "hello" {
		t.Errorf("default field: got %v, want %q", v, "hello")
	}
}
