package hclutil

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// CTYValueUnmarshaler is the interface implemented by types that can unmarshal
type CTYValueMarshaler interface {
	MarshalCTYValue() (cty.Value, error)
}

func MarshalCTYValue(v any) (cty.Value, error) {
	if m, ok := v.(CTYValueMarshaler); ok {
		return m.MarshalCTYValue()
	}
	rv := reflect.ValueOf(v)
	value, _, err := marshalCTYValue(rv)
	return value, err
}

var (
	marshalerType     = reflect.TypeOf((*CTYValueMarshaler)(nil)).Elem()
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func marshalCTYValue(rv reflect.Value) (cty.Value, bool, error) {
	if !rv.IsValid() {
		return cty.UnknownVal(cty.DynamicPseudoType), true, errors.New("invalid value")
	}
	rt := rv.Type()
	canAddr := rv.CanAddr()
	canInterface := rv.CanInterface()
	if rt.Kind() != reflect.Ptr && canAddr && reflect.PointerTo(rt).Implements(marshalerType) {
		m := rv.Addr().Interface().(CTYValueMarshaler)
		value, err := m.MarshalCTYValue()
		return value, value.IsNull() || !value.IsKnown(), err
	}
	if rt.Implements(marshalerType) && canInterface {
		m := rv.Interface().(CTYValueMarshaler)
		value, err := m.MarshalCTYValue()
		return value, value.IsNull() || !value.IsKnown(), err
	}
	if rt.Kind() != reflect.Ptr && canAddr && reflect.PtrTo(rt).Implements(jsonMarshalerType) {
		m := rv.Addr().Interface().(json.Marshaler)
		b, err := m.MarshalJSON()
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		val, err := jsonToCTYValue(b)
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		return val, val.IsNull(), nil
	}
	if rt.Implements(jsonMarshalerType) && canInterface {
		m := rv.Interface().(json.Marshaler)
		b, err := m.MarshalJSON()
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		val, err := jsonToCTYValue(b)
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		return val, false, nil
	}
	if rt.Kind() != reflect.Ptr && canAddr && reflect.PtrTo(rt).Implements(textMarshalerType) {
		m := rv.Addr().Interface().(encoding.TextMarshaler)
		b, err := m.MarshalText()
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		str := string(b)
		return cty.StringVal(str), str == "", nil
	}
	if rt.Implements(textMarshalerType) && canInterface {
		m := rv.Interface().(encoding.TextMarshaler)
		b, err := m.MarshalText()
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), false, err
		}
		return cty.StringVal(string(b)), false, nil
	}
	switch rv.Kind() {
	case reflect.Interface:
		if rv.IsNil() {
			return cty.NullVal(cty.DynamicPseudoType), true, nil
		}
		return marshalCTYValue(rv.Elem())
	case reflect.Ptr:
		if rv.IsNil() {
			switch rv.Type().Elem().Kind() {
			case reflect.String:
				return cty.NullVal(cty.String), true, nil
			case reflect.Bool:
				return cty.NullVal(cty.Bool), true, nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return cty.NullVal(cty.Number), true, nil
			default:
				return cty.NullVal(cty.DynamicPseudoType), true, nil
			}
		}
		return marshalCTYValue(rv.Elem())
	case reflect.Struct:
		field := getStructFileds(rv.Type())
		valueMap := make(map[string]cty.Value, len(field))
		for _, f := range field {
			fv := rv
			for _, i := range f.index {
				fv = fv.Field(i)
			}
			v, isEmpty, err := marshalCTYValue(fv)
			if err != nil {
				return cty.UnknownVal(cty.DynamicPseudoType), true, err
			}
			if isEmpty && f.omitEmpty {
				continue
			}
			valueMap[f.tagName] = v
		}
		return cty.ObjectVal(valueMap), len(valueMap) == 0, nil
	case reflect.Map:
		if rv.IsNil() {
			return cty.MapValEmpty(cty.DynamicPseudoType), true, nil
		}
		var keyHasMarshaler bool
		var keyHasPtrMarshaler bool
		if rt.Key().Kind() != reflect.String {
			if rt.Key().Kind() != reflect.Ptr && reflect.PtrTo(rt.Key()).Implements(marshalerType) {
				keyHasPtrMarshaler = true
			} else if rt.Key().Implements(marshalerType) {
				keyHasMarshaler = true
			} else {
				return cty.UnknownVal(cty.DynamicPseudoType), true, fmt.Errorf("unsupported map key type: %s", rt.Key())
			}
		}

		if rv.Len() == 0 {
			switch rt.Elem().Kind() {
			case reflect.String:
				return cty.MapValEmpty(cty.String), true, nil
			case reflect.Bool:
				return cty.MapValEmpty(cty.Bool), true, nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
				return cty.MapValEmpty(cty.Number), true, nil
			default:
				return cty.MapValEmpty(cty.DynamicPseudoType), true, nil
			}
		}
		valueMap := make(map[string]cty.Value, rv.Len())
		for _, key := range rv.MapKeys() {
			var keyStr string
			if keyHasMarshaler {
				m := key.Interface().(encoding.TextMarshaler)
				value, err := m.MarshalText()
				if err != nil {
					return cty.UnknownVal(cty.DynamicPseudoType), true, err
				}
				keyStr = string(value)
			} else if keyHasPtrMarshaler {
				if key.IsNil() || !key.CanAddr() {
					continue
				}
				m := key.Addr().Interface().(encoding.TextMarshaler)
				value, err := m.MarshalText()
				if err != nil {
					return cty.UnknownVal(cty.DynamicPseudoType), true, err
				}
				keyStr = string(value)
			} else {
				keyStr = key.String()
			}
			v, _, err := marshalCTYValue(rv.MapIndex(key))
			if err != nil {
				return cty.UnknownVal(cty.DynamicPseudoType), true, err
			}
			valueMap[keyStr] = v
		}
		return cty.MapVal(valueMap), len(valueMap) == 0, nil
	case reflect.Slice:
		return marshalCTYValueFromSlice(rv)
	case reflect.Array:
		return marshalCTYValueFromSlice(rv)
	case reflect.Bool:
		return cty.BoolVal(rv.Bool()), rv.IsZero(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return cty.NumberIntVal(rv.Int()), rv.IsZero(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return cty.NumberIntVal(int64(rv.Uint())), rv.IsZero(), nil
	case reflect.Float32, reflect.Float64:
		return cty.NumberFloatVal(rv.Float()), rv.IsZero(), nil
	case reflect.String:
		return cty.StringVal(rv.String()), rv.IsZero(), nil
	default:
		return cty.UnknownVal(cty.DynamicPseudoType), true, fmt.Errorf("unsupported type: %s", rt)
	}
}

func marshalCTYValueFromSlice(rv reflect.Value) (cty.Value, bool, error) {
	if rv.IsNil() {
		return cty.ListValEmpty(cty.DynamicPseudoType), true, nil
	}
	valueList := make([]cty.Value, rv.Len())
	elemType := cty.DynamicPseudoType
	elemCount := 0
	isTuple := false
	for i := 0; i < rv.Len(); i++ {
		elemCount++
		v, _, err := marshalCTYValue(rv.Index(i))
		if err != nil {
			return cty.UnknownVal(cty.DynamicPseudoType), true, err
		}
		valueList[i] = v
		if isTuple {
			continue
		}
		if !v.Type().IsPrimitiveType() {
			isTuple = true
			continue
		}
		if elemType == cty.DynamicPseudoType {
			elemType = v.Type()
		} else {
			if elemType != v.Type() {
				isTuple = true
			}
		}
	}
	if elemCount == 0 {
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			return cty.ListValEmpty(cty.String), true, nil
		case reflect.Bool:
			return cty.ListValEmpty(cty.Bool), true, nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			return cty.ListValEmpty(cty.Number), true, nil
		default:
			return cty.EmptyTupleVal, true, nil
		}
	}
	if isTuple {
		return cty.TupleVal(valueList), false, nil
	}
	return cty.ListVal(valueList), false, nil
}
