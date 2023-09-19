package hclutil_test

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/Songmu/flextime"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
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

func TestHCLFunctionCoalesce__MarshalNull(t *testing.T) {
	restore := flextime.Fix(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	defer restore()
	expr, diags := hclsyntax.ParseExpression([]byte(`coalesce(var.start_at,now())`), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		t.Log(diags)
		t.Fatal("parse failed")
	}
	ctx := hclutil.NewEvalContext()
	var num *int64
	val, err := hclutil.MarshalCTYValue(num)
	if err != nil {
		t.Log(err)
		t.Fatal("marshal failed")
	}
	ctx.Variables = map[string]cty.Value{
		"var": cty.ObjectVal(map[string]cty.Value{
			"start_at": val,
		}),
	}
	var ret int64
	diags = gohcl.DecodeExpression(expr, ctx, &ret)
	if diags.HasErrors() {
		t.Log(diags)
		t.Fatal("decode failed")
	}
	if ret != 1704067200 {
		t.Errorf("want %d, got %d", 1704067200, ret)
	}
}
