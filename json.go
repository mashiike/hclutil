package hclutil

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var (
	encodeExprStr = `jsonencode(var)`
	encodeExpr    hcl.Expression
	decodeExprStr = `jsondecode(var)`
	decodeExpr    hcl.Expression
)

func init() {
	var diags hcl.Diagnostics
	encodeExpr, diags = hclsyntax.ParseExpression([]byte(encodeExprStr), "cty_value_to_json.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		panic(fmt.Sprintf("hclutil: failed to parse expression %s: %s", encodeExprStr, diags.Error()))
	}
	decodeExpr, diags = hclsyntax.ParseExpression([]byte(decodeExprStr), "json_to_cty_value.hcl", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		panic(fmt.Sprintf("hclutil: failed to parse expression %s: %s", decodeExprStr, diags.Error()))
	}
}

func ctyValueToJSON(value cty.Value) ([]byte, error) {
	if !value.IsKnown() {
		return nil, &UnknownValueError{Value: value}
	}
	v, diags := encodeExpr.Value(&hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": value,
		},
		Functions: map[string]function.Function{
			"jsonencode": stdlib.JSONEncodeFunc,
		},
	})
	if diags.HasErrors() {
		return nil, fmt.Errorf("convert cty.Value to JSON: %w", diags)
	}
	return []byte(v.AsString()), nil
}

func jsonToCTYValue(bs []byte) (cty.Value, error) {
	v, diags := decodeExpr.Value(&hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.StringVal(string(bs)),
		},
		Functions: map[string]function.Function{
			"jsondecode": stdlib.JSONDecodeFunc,
		},
	})
	if diags.HasErrors() {
		return cty.NullVal(cty.DynamicPseudoType), fmt.Errorf("convert cty.Value to JSON: %w", diags)
	}
	return v, nil
}
