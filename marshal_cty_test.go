package hclutil_test

import (
	"testing"

	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestMarshalCTYValue__Primitive(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		v := "hoge"
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		if got != cty.StringVal("hoge") {
			t.Errorf("got = %s, want hoge", got)
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		v := 1234
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.NumberIntVal(1234)
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsBigFloat().String() != want.AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		v := 1234.5678
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.NumberFloatVal(1234.5678)
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsBigFloat().String() != want.AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		v := true
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		if got != cty.True {
			t.Errorf("got = %s, want true", got)
		}
	})
}

func TestMarshalCTYValue__PtrPrimitive(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		v := "hoge"
		got, err := hclutil.MarshalCTYValue(&v)
		if err != nil {
			t.Error(err)
		}
		if got != cty.StringVal("hoge") {
			t.Errorf("got = %s, want hoge", got)
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		v := 1234
		got, err := hclutil.MarshalCTYValue(&v)
		if err != nil {
			t.Error(err)
		}
		want := cty.NumberIntVal(1234)
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsBigFloat().String() != want.AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		v := 1234.5678
		got, err := hclutil.MarshalCTYValue(&v)
		if err != nil {
			t.Error(err)
		}
		want := cty.NumberFloatVal(1234.5678)
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsBigFloat().String() != want.AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		v := true
		got, err := hclutil.MarshalCTYValue(&v)
		if err != nil {
			t.Error(err)
		}
		if got != cty.True {
			t.Errorf("got = %s, want true", got)
		}
	})
}

func TestMarshalCTYValue__Slice(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		v := []string{"hoge", "fuga"}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.ListVal([]cty.Value{cty.StringVal("hoge"), cty.StringVal("fuga")})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueSlice()[0] != want.AsValueSlice()[0] {
			t.Errorf("got = %s, want %s", got, want)
		}
		if got.AsValueSlice()[1] != want.AsValueSlice()[1] {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		v := []int{1234, 5678}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.ListVal([]cty.Value{cty.NumberIntVal(1234), cty.NumberIntVal(5678)})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		sliceGot := got.AsValueSlice()
		sliceWant := want.AsValueSlice()
		if len(sliceGot) != len(sliceWant) {
			t.Errorf("got = %s, want %s", got, want)
		}
		if sliceGot[0].AsBigFloat().String() != sliceWant[0].AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
		if sliceGot[1].AsBigFloat().String() != sliceWant[1].AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		v := []float64{1234.5678, 5678.1234}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.ListVal([]cty.Value{
			cty.NumberFloatVal(1234.5678),
			cty.NumberFloatVal(5678.1234),
		})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		sliceGot := got.AsValueSlice()
		sliceWant := want.AsValueSlice()
		if len(sliceGot) != len(sliceWant) {
			t.Errorf("got = %s, want %s", got, want)
		}
		if sliceGot[0].AsBigFloat().String() != sliceWant[0].AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
		if sliceGot[1].AsBigFloat().String() != sliceWant[1].AsBigFloat().String() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		v := []bool{true, false}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.ListVal([]cty.Value{cty.True, cty.False})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueSlice()[0] != want.AsValueSlice()[0] {
			t.Errorf("got = %s, want %s", got, want)
		}
		if got.AsValueSlice()[1] != want.AsValueSlice()[1] {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("interface{}", func(t *testing.T) {
		t.Parallel()
		v := []interface{}{1234, "hoge", true}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.TupleVal([]cty.Value{
			cty.NumberIntVal(1234),
			cty.StringVal("hoge"),
			cty.True,
		})
		if got.Type().IsTupleType() != want.Type().IsTupleType() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
			t.FailNow()
		}
		sliceGot := got.AsValueSlice()
		sliceWant := want.AsValueSlice()
		if len(sliceGot) != len(sliceWant) {
			t.Errorf("got = %s, want %s", got, want)
			t.FailNow()
		}
		for i := range sliceGot {
			if sliceGot[i].Type() != sliceWant[i].Type() {
				t.Errorf("got = %s, want %s", got, want)
				t.FailNow()
			}
			if sliceGot[i].GoString() != sliceWant[i].GoString() {
				t.Errorf("got = %s, want %s", got, want)
				t.FailNow()
			}
		}
	})
}

func TestMarshalCTYValue__Map(t *testing.T) {
	t.Parallel()
	t.Run("string", func(t *testing.T) {
		t.Parallel()
		v := map[string]string{"hoge": "fuga"}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.MapVal(map[string]cty.Value{"hoge": cty.StringVal("fuga")})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueMap()["hoge"] != want.AsValueMap()["hoge"] {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		v := map[string]int{"hoge": 1234}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.MapVal(map[string]cty.Value{"hoge": cty.NumberIntVal(1234)})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueMap()["hoge"].GoString() != want.AsValueMap()["hoge"].GoString() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("float", func(t *testing.T) {
		t.Parallel()
		v := map[string]float64{"hoge": 1234.5678}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.MapVal(map[string]cty.Value{"hoge": cty.NumberFloatVal(1234.5678)})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueMap()["hoge"].GoString() != want.AsValueMap()["hoge"].GoString() {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		v := map[string]bool{"hoge": true}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
		}
		want := cty.MapVal(map[string]cty.Value{"hoge": cty.True})
		if got.Type() != want.Type() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
		}
		if got.AsValueMap()["hoge"] != want.AsValueMap()["hoge"] {
			t.Errorf("got = %s, want %s", got, want)
		}
	})
	t.Run("interface{}", func(t *testing.T) {
		t.Parallel()
		v := map[string]interface{}{"hoge": 1234}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.MapVal(map[string]cty.Value{"hoge": cty.NumberIntVal(1234)})
		if got.Type().IsMapType() != want.Type().IsMapType() {
			t.Errorf("got type = %s, want %s", got.Type(), want.Type())
			t.FailNow()
		}
		mapGot := got.AsValueMap()
		mapWant := want.AsValueMap()
		if len(mapGot) != len(mapWant) {
			t.Errorf("got = %s, want %s", got, want)
			t.FailNow()
		}
		for k := range mapGot {
			if mapGot[k].Type() != mapWant[k].Type() {
				t.Errorf("got = %s, want %s", got, want)
				t.FailNow()
			}
			if mapGot[k].GoString() != mapWant[k].GoString() {
				t.Errorf("got = %s, want %s", got, want)
				t.FailNow()
			}
		}
	})
}

func TestMarshalCTYValue__Struct(t *testing.T) {
	t.Parallel()
	t.Run("no tag", func(t *testing.T) {
		t.Parallel()
		type embedded struct {
			Embedded string
		}
		type hoge struct {
			Foo         string
			FooEmpty    string
			Bar         int
			BarEmpty    int
			FooBar      bool
			FooBarEmpty bool
			unexported  string
			embedded
		}
		v := hoge{
			Foo:        "fuga",
			Bar:        1234,
			FooBar:     false,
			unexported: "unexported",
			embedded: embedded{
				Embedded: "embedded",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.ObjectVal(map[string]cty.Value{
			"foo":           cty.StringVal("fuga"),
			"foo_empty":     cty.StringVal(""),
			"bar":           cty.NumberIntVal(1234),
			"bar_empty":     cty.NumberIntVal(0),
			"foo_bar":       cty.BoolVal(false),
			"foo_bar_empty": cty.BoolVal(false),
			"embedded":      cty.StringVal("embedded"),
		})
		t.Log(got.GoString())
		if !got.Type().IsObjectType() {
			t.Errorf("got is not object type: %s", got.Type().GoString())
			t.FailNow()
		}
		gotMap := got.AsValueMap()
		wantMap := want.AsValueMap()
		if len(gotMap) != len(wantMap) {
			t.Errorf("got length = %d, want %d", len(gotMap), len(wantMap))
			t.FailNow()
		}
		for k := range gotMap {
			if gotMap[k].Type().GoString() != wantMap[k].Type().GoString() {
				t.Errorf("got type = %s, want %s", gotMap[k].Type(), wantMap[k].Type())
				t.FailNow()
			}
			if gotMap[k].GoString() != wantMap[k].GoString() {
				t.Errorf("got = %s, want %s", gotMap[k], wantMap[k])
				t.FailNow()
			}
		}
	})
	t.Run("tag", func(t *testing.T) {
		t.Parallel()
		type embedded struct {
			Embedded string `cty:"embedded"`
			Ignored  string `cty:"-"`
		}
		type hoge struct {
			Foo         string `cty:"foo"`
			FooEmpty    string `cty:"foo_empty,omitempty"`
			Bar         int    `cty:"bar"`
			BarEmpty    int    `cty:"bar_empty,omitempty"`
			FooBar      bool   `cty:"foo_bar"`
			FooBarEmpty bool   `cty:"foo_bar_empty,omitempty"`
			unexported  string `cty:"unexported,omitempty"`
			IgnoredStr  string `cty:"-"`
			embedded
		}
		v := hoge{
			Foo:        "fuga",
			Bar:        1234,
			FooBar:     false,
			unexported: "unexported",
			IgnoredStr: "ignored",
			embedded: embedded{
				Embedded: "embedded",
				Ignored:  "ignored",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.ObjectVal(map[string]cty.Value{
			"foo":      cty.StringVal("fuga"),
			"bar":      cty.NumberIntVal(1234),
			"foo_bar":  cty.BoolVal(false),
			"embedded": cty.StringVal("embedded"),
		})
		if !got.Type().IsObjectType() {
			t.Errorf("got is not object type: %s", got.Type().GoString())
			t.FailNow()
		}
		gotMap := got.AsValueMap()
		wantMap := want.AsValueMap()
		if len(gotMap) != len(wantMap) {
			t.Errorf("got length = %d, want %d", len(gotMap), len(wantMap))
			t.FailNow()
		}
		for k := range gotMap {
			if gotMap[k].Type().GoString() != wantMap[k].Type().GoString() {
				t.Errorf("got type = %v, want %v", gotMap[k].Type(), wantMap[k].Type())
				t.FailNow()
			}
			if gotMap[k].GoString() != wantMap[k].GoString() {
				t.Errorf("got = %s, want %s", gotMap[k], wantMap[k])
				t.FailNow()
			}
		}
	})
}

type testCTYValueMarshaler struct {
	Val string
}

func (t testCTYValueMarshaler) MarshalCTYValue() (cty.Value, error) {
	return cty.StringVal(t.Val), nil
}

func TestMarshalCTYValue__Marshaler(t *testing.T) {
	t.Parallel()
	t.Run("struct", func(t *testing.T) {
		v := testCTYValueMarshaler{
			Val: "hoge",
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.StringVal("hoge")
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("pointer", func(t *testing.T) {
		v := &testCTYValueMarshaler{
			Val: "hoge",
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.StringVal("hoge")
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("slice", func(t *testing.T) {
		v := []testCTYValueMarshaler{
			{
				Val: "hoge",
			},
			{
				Val: "fuga",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.ListVal([]cty.Value{
			cty.StringVal("hoge"),
			cty.StringVal("fuga"),
		})
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("map", func(t *testing.T) {
		v := map[string]testCTYValueMarshaler{
			"hoge": {
				Val: "fuga",
			},
			"foo": {
				Val: "bar",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.MapVal(map[string]cty.Value{
			"hoge": cty.StringVal("fuga"),
			"foo":  cty.StringVal("bar"),
		})
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
}

type testTextMarshaler struct {
	Val string
}

func (t testTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(t.Val), nil
}

func TestMarshalCTYValue__TextMarshaler(t *testing.T) {
	t.Parallel()
	t.Run("struct", func(t *testing.T) {
		v := testTextMarshaler{
			Val: "hoge",
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.StringVal("hoge")
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("pointer", func(t *testing.T) {
		v := &testTextMarshaler{
			Val: "hoge",
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.StringVal("hoge")
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("slice", func(t *testing.T) {
		v := []testTextMarshaler{
			{
				Val: "hoge",
			},
			{
				Val: "fuga",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.ListVal([]cty.Value{
			cty.StringVal("hoge"),
			cty.StringVal("fuga"),
		})
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
	t.Run("map", func(t *testing.T) {
		v := map[string]testTextMarshaler{
			"hoge": {
				Val: "fuga",
			},
			"foo": {
				Val: "bar",
			},
		}
		got, err := hclutil.MarshalCTYValue(v)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		want := cty.MapVal(map[string]cty.Value{
			"hoge": cty.StringVal("fuga"),
			"foo":  cty.StringVal("bar"),
		})
		if got.GoString() != want.GoString() {
			t.Errorf("got = %s, want %s", got.GoString(), want.GoString())
			t.FailNow()
		}
	})
}
