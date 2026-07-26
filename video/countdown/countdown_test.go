package countdown

import (
	"testing"
	"time"
)

func TestParseTarget(t *testing.T) {

	want := time.Date(2027, time.January, 1, 0, 0, 0, 0, jst)

	cases := []struct {
		name string
		in   string
	}{
		{"rfc3339", "2027-01-01T00:00:00+09:00"},
		// ikasbox の登録 UI の <input type="datetime-local"> が出す形
		{"datetime-local", "2027-01-01T00:00"},
		{"datetime-local with seconds", "2027-01-01T00:00:00"},
		{"space separated", "2027-01-01 00:00:00"},
		{"space separated without seconds", "2027-01-01 00:00"},
		{"surrounding spaces", "  2027-01-01T00:00  "},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseTarget(c.in)
			if err != nil {
				t.Fatalf("parseTarget(%q) error: %v", c.in, err)
			}
			if !got.Equal(want) {
				t.Errorf("parseTarget(%q) = %v, want %v", c.in, got, want)
			}
		})
	}
}

func TestParseTargetInvalid(t *testing.T) {
	for _, in := range []string{"not-a-date", "2027-13-01T00:00", ""} {
		if _, err := parseTarget(in); err == nil {
			t.Errorf("parseTarget(%q): expected error", in)
		}
	}
}

// New は target 未指定を許容する(過去日時扱いで即終了表示)
func TestNewWithoutTarget(t *testing.T) {
	v, err := New(`{"text":"done"}`)
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	defer v.Release()

	if v.text != "done" {
		t.Errorf("text = %q, want %q", v.text, "done")
	}
}

func TestNewInvalidTarget(t *testing.T) {
	if _, err := New(`{"target":"not-a-date"}`); err == nil {
		t.Fatal("expected error for invalid target")
	}
}

func TestSpecFields(t *testing.T) {
	fields := Spec()
	if len(fields) == 0 {
		t.Fatal("Spec() returned no fields")
	}

	byName := map[string]string{}
	for _, f := range fields {
		byName[f.Name] = string(f.Type)
	}

	if byName["target"] != "datetime" {
		t.Errorf("target field type = %q, want %q", byName["target"], "datetime")
	}
	if _, ok := byName["text"]; !ok {
		t.Errorf("Spec() missing 'text' field: %+v", fields)
	}
}
