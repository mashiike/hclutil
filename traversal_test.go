package hclutil_test

import (
	"testing"

	"github.com/mashiike/hclutil"
)

func TestVariablesReffarances(t *testing.T) {
	src := []byte(`
[
	hoge.fuga,
	hoge.fuga[0],
	hoge.fuga[0].piyo,
	hoge.fuga["piyo"].piyo,
]`)
	expr, diags := hclutil.ParseExpression(src)
	if diags.HasErrors() {
		t.Fatal(diags)
	}
	vars := hclutil.VariablesReffarances(expr)
	if len(vars) != 4 {
		t.Fatalf("unexpected variables length: %d", len(vars))
	}
	if vars[0] != "hoge.fuga" {
		t.Errorf("unexpected variable: %s", vars[0])
	}
	if vars[1] != "hoge.fuga[0]" {
		t.Errorf("unexpected variable: %s", vars[1])
	}
	if vars[2] != "hoge.fuga[0].piyo" {
		t.Errorf("unexpected variable: %s", vars[2])
	}
	if vars[3] != "hoge.fuga[\"piyo\"].piyo" {
		t.Errorf("unexpected variable: %s", vars[3])
	}
}
