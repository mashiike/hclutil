package hclutil_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestUnmarshalCTYValue__PrimitiveType(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		var v string
		err := hclutil.UnmarshalCTYValue(cty.StringVal("hoge"), &v)
		if err != nil {
			t.Error(err)
		}
		if v != "hoge" {
			t.Errorf("v = %s, want hoge", v)
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		var v int
		err := hclutil.UnmarshalCTYValue(cty.NumberIntVal(1234), &v)
		if err != nil {
			t.Error(err)
		}
		if v != 1234 {
			t.Errorf("v = %d, want 1234", v)
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		var v float64
		err := hclutil.UnmarshalCTYValue(cty.NumberFloatVal(1234.5678), &v)
		if err != nil {
			t.Error(err)
		}
		if v != 1234.5678 {
			t.Errorf("v = %f, want 1234.5678", v)
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		var v bool
		err := hclutil.UnmarshalCTYValue(cty.True, &v)
		if err != nil {
			t.Error(err)
		}
		if !v {
			t.Errorf("v = %t, want true", v)
		}
	})
}

func TestUnmarshalCTYValue__PtrPrimitive(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		var v *string
		err := hclutil.UnmarshalCTYValue(cty.StringVal("hoge"), &v)
		if err != nil {
			t.Error(err)
		}
		if v == nil {
			t.Error("v is nil")
		} else {
			if *v != "hoge" {
				t.Errorf("*v = %s, want hoge", *v)
			}
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		var v *int
		err := hclutil.UnmarshalCTYValue(cty.NumberIntVal(1234), &v)
		if err != nil {
			t.Error(err)
		}
		if v == nil {
			t.Error("v is nil")
		} else {
			if *v != 1234 {
				t.Errorf("*v = %d, want 1234", *v)
			}
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		var v *float64
		err := hclutil.UnmarshalCTYValue(cty.NumberFloatVal(1234.5678), &v)
		if err != nil {
			t.Error(err)
		}
		if v == nil {
			t.Error("v is nil")
		} else {
			if *v != 1234.5678 {
				t.Errorf("*v = %f, want 1234.5678", *v)
			}
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		var v *bool
		err := hclutil.UnmarshalCTYValue(cty.True, &v)
		if err != nil {
			t.Error(err)
		}
		if v == nil {
			t.Error("v is nil")
		} else {
			if !*v {
				t.Errorf("*v = %t, want true", *v)
			}
		}
	})
}

func TestUnmarshalCTYValue__Interface(t *testing.T) {
	t.Parallel()
	t.Run("string to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(cty.StringVal("hoge"), &v)
		if err != nil {
			t.Error(err)
		}
		if v != "hoge" {
			t.Errorf("v = %s, want hoge", v)
		}
	})
	t.Run("integer to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(cty.NumberIntVal(1234), &v)
		if err != nil {
			t.Error(err)
		}
		bigFloat, ok := v.(*big.Float)
		if !ok {
			t.Errorf("v is not *big.Float")
		}
		if bigFloat.Cmp(big.NewFloat(1234)) != 0 {
			t.Errorf("v = %s, want 1234", bigFloat)
		}
	})
	t.Run("float to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(cty.NumberFloatVal(1234.5678), &v)
		if err != nil {
			t.Error(err)
		}
		bigFloat, ok := v.(*big.Float)
		if !ok {
			t.Errorf("v is not *big.Float")
		}
		if bigFloat.Cmp(big.NewFloat(1234.5678)) != 0 {
			t.Errorf("v = %s, want 1234.5678", bigFloat)
		}
	})
	t.Run("bool to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(cty.True, &v)
		if err != nil {
			t.Error(err)
		}
		if v != true {
			t.Errorf("v = %t, want true", v)
		}
	})
	t.Run("list to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(
			cty.ListVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		vList, ok := v.([]interface{})
		if !ok {
			t.Errorf("v is not []interface{}")
		}
		if len(vList) != 2 {
			t.Errorf("len(vList) = %d, want 2", len(vList))
		}
		if vList[0] != "foo" {
			t.Errorf("vList[0] = %s, want foo", vList[0])
		}
		if vList[1] != "bar" {
			t.Errorf("vList[1] = %s, want bar", vList[1])
		}
	})
	t.Run("map to interface{}", func(t *testing.T) {
		t.Parallel()
		var v interface{}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		vMap, ok := v.(map[string]interface{})
		if !ok {
			t.Errorf("v is not map[string]interface{}")
		}
		if vMap["foo"] != "bar" {
			t.Errorf("vMap[\"foo\"] = %s, want bar", vMap["foo"])
		}
	})
	t.Run("map[string]interface{}", func(t *testing.T) {
		t.Parallel()
		var v map[string]interface{}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v["foo"] != "bar" {
			t.Errorf("v[\"foo\"] = %s, want bar", v["foo"])
		}
	})
	t.Run("[]interface{}", func(t *testing.T) {
		t.Parallel()
		var v []interface{}
		err := hclutil.UnmarshalCTYValue(
			cty.ListVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if len(v) != 2 {
			t.Errorf("len(v) = %d, want 2", len(v))
		}
		if v[0] != "foo" {
			t.Errorf("v[0] = %s, want foo", v[0])
		}
		if v[1] != "bar" {
			t.Errorf("v[1] = %s, want bar", v[1])
		}
	})
	t.Run("[][]interface{}", func(t *testing.T) {
		t.Parallel()
		var v [][]interface{}
		err := hclutil.UnmarshalCTYValue(
			cty.TupleVal([]cty.Value{
				cty.TupleVal([]cty.Value{
					cty.StringVal("foo"),
					cty.NumberIntVal(1234),
				}),
				cty.TupleVal([]cty.Value{
					cty.StringVal("baz"),
					cty.NumberIntVal(5678),
				}),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if len(v) != 2 {
			t.Errorf("len(v) = %d, want 2", len(v))
		}
		if len(v[0]) != 2 {
			t.Errorf("len(v[0]) = %d, want 2", len(v[0]))
		}
		if v[0][0] != "foo" {
			t.Errorf("v[0][0] = %s, want foo", v[0][0])
		}
		if f, ok := v[0][1].(*big.Float); ok && f.String() != big.NewInt(1234).String() {
			t.Errorf("v[0][1] = %v, want 1234", v[0][1])
		}
		if len(v[1]) != 2 {
			t.Errorf("len(v[1]) = %d, want 2", len(v[1]))
		}
		if v[1][0] != "baz" {
			t.Errorf("v[1][0] = %s, want baz", v[1][0])
		}
		if f, ok := v[1][1].(*big.Float); ok && f.String() != big.NewInt(5678).String() {
			t.Errorf("v[1][1] = %v, want 5678", v[1][1])
		}
	})
}

func TestUnmarshalCTYValue__Struct(t *testing.T) {
	t.Parallel()
	t.Run("no tag struct", func(t *testing.T) {
		t.Parallel()
		var v struct {
			Foo    string
			Bar    int
			FooBar bool
		}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo":     cty.StringVal("bar"),
				"bar":     cty.NumberIntVal(1234),
				"foo_bar": cty.BoolVal(true),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Foo != "bar" {
			t.Errorf("v.Foo = %s, want bar", v.Foo)
		}
		if v.Bar != 1234 {
			t.Errorf("v.Bar = %d, want 1234", v.Bar)
		}
		if v.FooBar != true {
			t.Errorf("v.FooBar = %t, want true", v.FooBar)
		}
	})
	t.Run("tag struct", func(t *testing.T) {
		t.Parallel()
		type embedded struct {
			Embeded string `cty:"embeded"`
		}
		var v struct {
			embedded
			Foo        string `cty:"foo"`
			Bar        int    `cty:"bar"`
			FooBar     bool   `cty:"baza"`
			Ignore     string `cty:"-"`
			unexported string
		}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo":     cty.StringVal("bar"),
				"bar":     cty.NumberIntVal(1234),
				"baza":    cty.BoolVal(true),
				"embeded": cty.StringVal("embeded"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Foo != "bar" {
			t.Errorf("v.Foo = %s, want bar", v.Foo)
		}
		if v.Bar != 1234 {
			t.Errorf("v.Bar = %d, want 1234", v.Bar)
		}
		if v.FooBar != true {
			t.Errorf("v.FooBar = %t, want true", v.FooBar)
		}
		if v.Ignore != "" {
			t.Errorf("v.Ignore = %s, want \"\"", v.Ignore)
		}
		if v.unexported != "" {
			t.Errorf("v.unexported = %s, want \"\"", v.unexported)
		}
		if v.Embeded != "embeded" {
			t.Errorf("v.Embeded = %s, want embeded", v.Embeded)
		}
	})
}

type testCTYValueUnmarshaler struct {
	Val string
}

func (t *testCTYValueUnmarshaler) UnmarshalCTYValue(v cty.Value) error {
	if v.Type() != cty.String {
		return fmt.Errorf("v.Type() = %s, want string", v.Type())
	}
	t.Val = v.AsString()
	return nil
}

func TestUnmarshalCTYValue__CTYValueUnmarshaler(t *testing.T) {
	t.Parallel()
	t.Run("primitive", func(t *testing.T) {
		t.Parallel()
		var v testCTYValueUnmarshaler
		err := hclutil.UnmarshalCTYValue(
			cty.StringVal("foo"),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Val != "foo" {
			t.Errorf("v.Val = %s, want foo", v.Val)
		}
	})
	t.Run("slice", func(t *testing.T) {
		t.Parallel()
		var v []testCTYValueUnmarshaler
		err := hclutil.UnmarshalCTYValue(
			cty.ListVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v[0].Val != "foo" {
			t.Errorf("v[0].Val = %s, want foo", v[0].Val)
		}
		if v[1].Val != "bar" {
			t.Errorf("v[1].Val = %s, want bar", v[1].Val)
		}
	})
	t.Run("struct", func(t *testing.T) {
		t.Parallel()
		var v struct {
			Foo testCTYValueUnmarshaler `cty:"foo"`
		}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Foo.Val != "bar" {
			t.Errorf("v.Foo.Val = %s, want bar", v.Foo.Val)
		}
	})
}

func TestUnmarshalCTYValue__JSONUnmarshaler(t *testing.T) {
	t.Parallel()
	t.Run("json.RawMessage", func(t *testing.T) {
		t.Parallel()
		var v json.RawMessage
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if string(v) != `{"foo":"bar"}` {
			t.Errorf("v = %s, want {\"foo\":\"bar\"}", v)
		}
	})
}

type testTextUnmarshaler struct {
	Val string
}

func (t *testTextUnmarshaler) UnmarshalText(text []byte) error {
	t.Val = string(text)
	return nil
}

func TestUnmarshalCTYValue__TextUnmarshaler(t *testing.T) {
	t.Parallel()
	t.Run("primitive", func(t *testing.T) {
		t.Parallel()
		var v testTextUnmarshaler
		err := hclutil.UnmarshalCTYValue(
			cty.StringVal("foo"),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Val != "foo" {
			t.Errorf("v.Val = %s, want foo", v.Val)
		}
	})
	t.Run("slice", func(t *testing.T) {
		t.Parallel()
		var v []testTextUnmarshaler
		err := hclutil.UnmarshalCTYValue(
			cty.ListVal([]cty.Value{
				cty.StringVal("foo"),
				cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v[0].Val != "foo" {
			t.Errorf("v[0].Val = %s, want foo", v[0].Val)
		}
		if v[1].Val != "bar" {
			t.Errorf("v[1].Val = %s, want bar", v[1].Val)
		}
	})
	t.Run("struct", func(t *testing.T) {
		t.Parallel()
		var v struct {
			Foo testTextUnmarshaler `cty:"foo"`
		}
		err := hclutil.UnmarshalCTYValue(
			cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			&v,
		)
		if err != nil {
			t.Error(err)
		}
		if v.Foo.Val != "bar" {
			t.Errorf("v.Foo.Val = %s, want bar", v.Foo.Val)
		}
	})
}
