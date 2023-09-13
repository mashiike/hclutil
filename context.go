package hclutil

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// NewEvalContext は よく使う基本的な関数を登録したEvalContextを作成します。
// NewEvalContext creates an EvalContext with basic functions.
func NewEvalContext(optFns ...func(*utilFunctionOptions)) *hcl.EvalContext {
	evalCtx := &hcl.EvalContext{}
	evalCtx = WithUtilFunctions(evalCtx, optFns...)
	return evalCtx
}

// WithVariables returns a new EvalContext with parent's variables and variables merged.
func WithVariables(ctx *hcl.EvalContext, variables map[string]cty.Value) *hcl.EvalContext {
	cctx := ctx.NewChild()
	variablesList := make([]map[string]cty.Value, 0)
	for current := ctx; current != nil; current = current.Parent() {
		if current.Variables == nil {
			continue
		}
		variablesList = append([]map[string]cty.Value{current.Variables}, variablesList...)
	}
	variablesList = append(variablesList, variables)
	cctx.Variables = MergeVariables(variablesList...)
	return cctx
}

// WithValue returns a new EvalContext with path's value set.
// ctx = WithValue(ctx, "a.b.c", cty.StringVal("hoge"))
func WithValue(ctx *hcl.EvalContext, path string, value cty.Value) *hcl.EvalContext {
	paths := strings.Split(path, ".")
	variables := make(map[string]cty.Value, 1)
	variables[paths[len(paths)-1]] = value
	for i := len(paths) - 2; i >= 0; i-- {
		variables = map[string]cty.Value{paths[i]: cty.ObjectVal(variables)}
	}
	return WithVariables(ctx, variables)
}
