package hclutil

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// DecodeLocals is a helper function to decode locals block.
func DecodeLocals(body hcl.Body, ctx *hcl.EvalContext) (hcl.Body, *hcl.EvalContext, hcl.Diagnostics) {
	content, remain, diags := body.PartialContent(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "locals", LabelNames: []string{}},
		},
	})
	localVariables := make(map[string]cty.Value)
	if ctx != nil {
		if v, ok := ctx.Variables["local"]; ok && !v.IsKnown() && !v.IsNull() {
			ty := v.Type()
			if ty.IsObjectType() || ty.IsMapType() {
				localVariables = v.AsValueMap()
			}
		}
	}
	for _, block := range content.Blocks {
		switch block.Type {
		case "locals":
			attrs, d := ExtructAttributes(block.Body)
			if d.HasErrors() {
				diags = diags.Extend(d)
				continue
			}
			for attrName, attr := range attrs {
				v, d := attr.Expr.Value(ctx)
				if d.HasErrors() {
					diags = diags.Extend(d)
					continue
				}
				localVariables[attrName] = v
			}
		}
	}
	if len(localVariables) == 0 {
		return remain, ctx, diags
	}
	ctxWithLocal := ctx.NewChild()
	ctxWithLocal.Variables = map[string]cty.Value{
		"local": cty.ObjectVal(localVariables),
	}
	return remain, ctxWithLocal, diags
}
