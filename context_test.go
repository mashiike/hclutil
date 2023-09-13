package hclutil_test

import (
	"reflect"
	"testing"

	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestWithValue(t *testing.T) {
	t.Parallel()

	ctx := hclutil.NewEvalContext()
	ctx = hclutil.WithValue(ctx, "a.b.c", cty.StringVal("hoge"))

	got := ctx.Variables
	want := map[string]cty.Value{
		"a": cty.ObjectVal(map[string]cty.Value{
			"b": cty.ObjectVal(map[string]cty.Value{
				"c": cty.StringVal("hoge"),
			}),
		}),
	}
	var actual, expectd map[string]interface{}
	if err := hclutil.UnmarshalCTYValue(cty.ObjectVal(got), &actual); err != nil {
		t.Errorf("got error: %s", err)
	}
	if err := hclutil.UnmarshalCTYValue(cty.ObjectVal(want), &expectd); err != nil {
		t.Errorf("got error: %s", err)
	}
	if !reflect.DeepEqual(actual, expectd) {
		t.Errorf("got: %#v, want: %#v", actual, expectd)
	}
}
