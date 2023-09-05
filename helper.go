package hclutil

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// DumpCTYValue は cty.Value をJSON文字列に変換します。
//
//	これは、ログ出力等を行うときのデバッグ用途を想定しています。
func DumpCTYValue(v cty.Value) (string, error) {
	expr, diags := hclsyntax.ParseExpression([]byte(`jsonencode(var)`), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return "", diags
	}
	evalCtx := NewEvalContext()
	evalCtx.Variables = map[string]cty.Value{
		"var": v,
	}
	result, diags := expr.Value(evalCtx)
	if diags.HasErrors() {
		return "", diags
	}
	return result.AsString(), nil
}

// MustDumpCtyValue は cty.Value をJSON文字列に変換します。
//
// DumpCTYValue と異なり、エラーが発生した場合は panic します
func MustDumpCtyValue(v cty.Value) string {
	s, err := DumpCTYValue(v)
	if err != nil {
		panic(err)
	}
	return s
}
