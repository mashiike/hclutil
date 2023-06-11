package hclutil

import (
	"encoding"
	"reflect"
)

// このコードは、encoding/json/decode.goを参考にしています。
// オリジナルのコードのライセンスは以下の通りです。
// 　　　https://github.com/golang/go/blob/master/LICENSE
// コードの改変内容は、json.UnmarshalerをCTYValueUnmarshalerに変更したこととです。
//
// indirect は、必要に応じてポインタを割り当てながら v を下に移動し、ポインタでないものになるまで移動します。
// その過程で、CTYUnmarshalerに遭遇したらその時点で打ち切ります。
// また、decodingTextが有効である場合は、encoding.TextUnmarshalerを実装している場合にも打ち切ります。j@w
// decodingNullがtrueの場合は、nilを設定できる最初のポインタで打ち切ります。
func indirect(v reflect.Value, decodingNull bool) (CTYValueUnmarshaler, encoding.TextUnmarshaler, reflect.Value) {
	v0 := v
	haveAddr := false

	// もしvが名前付き型でアドレスを持っているなら、ポインタメソッドを持っている場合に見つけるために、アドレスから始めます。
	if v.Kind() != reflect.Pointer && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Pointer && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Pointer) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Pointer {
			break
		}

		if decodingNull && v.CanSet() {
			break
		}

		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(CTYValueUnmarshaler); ok {
				return u, nil, reflect.Value{}
			}
			if !decodingNull {
				if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
					return nil, u, reflect.Value{}
				}
			}
		}

		if haveAddr {
			v = v0
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, nil, v
}
