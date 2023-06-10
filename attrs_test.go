package hclutil_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mashiike/hclutil"
)

func TestExtructAttributes(t *testing.T) {
	t.Parallel()
	src := `
text = "hoge"
value = 1234
yesno = true

section "hoge" {
	text = "hoge"
	depends_on = [
		section.fuga,
	]
}

section "fuga" {
	text = "fuga"
}`
	file, _ := hclsyntax.ParseConfig([]byte(src), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	attrs, diags := hclutil.ExtructAttributes(file.Body)
	diagsReport(t, diags)
	if len(attrs) != 3 {
		t.Errorf("len(attrs) = %d, want 3", len(attrs))
	}
	if _, ok := attrs["text"]; !ok {
		t.Errorf("attrs does not have text")
	}
	if _, ok := attrs["value"]; !ok {
		t.Errorf("attrs does not have value")
	}
	if _, ok := attrs["yesno"]; !ok {
		t.Errorf("attrs does not have yesno")
	}
}

func TestExtructAttributes__MergedBody(t *testing.T) {
	t.Parallel()
	src := `
text = "hoge"

section "hoge" {
	text = "hoge"
	depends_on = [
		section.fuga,
	]
}`
	file, _ := hclsyntax.ParseConfig([]byte(src), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	src2 := `
value = 1234
yesno = true

section "hoge" {
	text = "hoge"
	depends_on = [
		section.fuga,
	]
}

section "fuga" {
	text = "fuga"
}`
	file2, _ := hclsyntax.ParseConfig([]byte(src2), "test2.hcl", hcl.Pos{Line: 1, Column: 1})
	merged := hcl.MergeFiles([]*hcl.File{file, file2})
	attrs, diags := hclutil.ExtructAttributes(merged)
	diagsReport(t, diags)
	if len(attrs) != 3 {
		t.Errorf("len(attrs) = %d, want 3", len(attrs))
	}
	if _, ok := attrs["text"]; !ok {
		t.Errorf("attrs does not have text")
	}
	if _, ok := attrs["value"]; !ok {
		t.Errorf("attrs does not have value")
	}
	if _, ok := attrs["yesno"]; !ok {
		t.Errorf("attrs does not have yesno")
	}
}
