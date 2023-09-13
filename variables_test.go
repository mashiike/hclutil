package hclutil_test

import (
	"reflect"
	"testing"

	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestMergeVariables(t *testing.T) {
	got := hclutil.MergeVariables(
		map[string]cty.Value{
			"hoge": cty.StringVal("fuga"),
			"foo":  cty.StringVal("bar"),
			"baz": cty.ObjectVal(map[string]cty.Value{
				"qux":   cty.StringVal("quux"),
				"corge": cty.StringVal("grault"),
			}),
		},
		map[string]cty.Value{
			"hoge": cty.StringVal("fuga"),
			"baz": cty.ObjectVal(map[string]cty.Value{
				"piyo":  cty.StringVal("hogera"),
				"corge": cty.StringVal("yabai"),
			}),
		},
	)
	want := map[string]cty.Value{
		"hoge": cty.StringVal("fuga"),
		"foo":  cty.StringVal("bar"),
		"baz": cty.ObjectVal(map[string]cty.Value{
			"qux":   cty.StringVal("quux"),
			"piyo":  cty.StringVal("hogera"),
			"corge": cty.StringVal("yabai"),
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
