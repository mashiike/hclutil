package hclutil_test

import (
	"testing"
	"testing/fstest"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mashiike/hclutil"
)

func TestHCLFunctionFile__DutyPath(t *testing.T) {
	t.Parallel()
	testFs := fstest.MapFS{
		"hoge.txt": {
			Data: []byte(`hoge`),
			Mode: 0,
		},
	}
	var str string
	expr, diags := hclsyntax.ParseExpression([]byte(`file("./hoge/.././hoge.txt")`), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Log(diags)
		t.Fatal("parse failed")
	}
	diags = gohcl.DecodeExpression(expr, hclutil.NewEvalContext(
		hclutil.WithFS(testFs),
	), &str)
	if diags.HasErrors() {
		t.Log(diags)
		t.Fatal("decode failed")
	}
	if str != "hoge" {
		t.Errorf("want %q, got %q", "hoge", str)
	}
}
