package hclutil_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestDecodeLocals(t *testing.T) {
	t.Parallel()

	src := `
locals {
	hoge = "hoge"
	fuga = 1234
}

locals {
	piyo = {
		tora = "tora"
	}
}

text = local.hoge
value = local.fuga
piyotora = local.piyo.tora
`
	file, _ := hclsyntax.ParseConfig([]byte(src), "test.hcl", hcl.Pos{Line: 1, Column: 1})
	remain, ctx, diags := hclutil.DecodeLocals(file.Body, nil)
	diagsReport(t, diags)
	attrs, diags := hclutil.ExtructAttributes(remain)
	diagsReport(t, diags)
	if len(ctx.Variables) != 1 {
		t.Errorf("unexpected length: %d", len(ctx.Variables))
	}
	localValiable := ctx.Variables["local"]
	want := cty.ObjectVal(map[string]cty.Value{
		"hoge": cty.StringVal("hoge"),
		"fuga": cty.NumberIntVal(1234),
		"piyo": cty.ObjectVal(map[string]cty.Value{
			"tora": cty.StringVal("tora"),
		}),
	})
	if localValiable.GoString() != want.GoString() {
		t.Errorf("unexpected value: %s", localValiable.GoString())
	}
	if len(attrs) != 3 {
		t.Errorf("unexpected length: %d", len(attrs))
	}
}
