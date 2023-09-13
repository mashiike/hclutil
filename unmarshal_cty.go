package hclutil

import (
	"fmt"
	"reflect"

	"github.com/zclconf/go-cty/cty"
)

// UnmarshalCTYValue decodes a cty.Value into the value pointed to by v.
func UnmarshalCTYValue(value cty.Value, v any) error {
	if !value.IsKnown() {
		return &UnknownValueError{Value: value}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}
	return unmarshalCTYValue(value, rv)
}

func unmarshalCTYValue(value cty.Value, rv reflect.Value) error {
	t := value.Type()
	switch {
	case t.IsListType() || t.IsTupleType() || t.IsSetType():
		if rv.IsValid() {
			if err := unmarshalCTYList(value, rv); err != nil {
				return err
			}
		}
	case t.IsMapType() || t.IsObjectType():
		if rv.IsValid() {
			if err := unmarshalCTYObject(value, rv); err != nil {
				return err
			}
		}
	case t.IsPrimitiveType():
		if rv.IsValid() {
			if err := unmarshalCTYPrimitive(value, rv); err != nil {
				return err
			}
		}
	default:
		return &UnmarshalTypeError{CTYType: t, Type: rv.Type()}
	}
	return nil
}

func unmarshalCTYList(value cty.Value, rv reflect.Value) error {
	u, uj, ut, pv := indirect(rv, value.IsNull())
	if u != nil {
		return u.UnmarshalCTYValue(value)
	}
	if uj != nil {
		bs, err := ctyValueToJSON(value)
		if err != nil {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type(), Detail: err}
		}
		return uj.UnmarshalJSON(bs)
	}
	if ut != nil {
		return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type()}
	}

	switch pv.Kind() {
	case reflect.Interface:
		if pv.NumMethod() == 0 {
			converted, err := convertCTYList(value)
			if err != nil {
				return err
			}
			pv.Set(reflect.ValueOf(converted))
			return nil
		}
	case reflect.Array, reflect.Slice:
		valueSlice := value.AsValueSlice()
		if pv.Kind() == reflect.Slice {
			if pv.Cap() < len(valueSlice) {
				pv.Set(reflect.MakeSlice(pv.Type(), len(valueSlice), len(valueSlice)))
			}
			if pv.Len() < len(valueSlice) {
				pv.SetLen(len(valueSlice))
			}
		}
		for i, v := range valueSlice {
			if err := unmarshalCTYValue(v, pv.Index(i)); err != nil {
				return err
			}
		}
	default:
		return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
	}
	return nil
}

func unmarshalCTYObject(value cty.Value, rv reflect.Value) error {
	u, uj, ut, pv := indirect(rv, value.IsNull())
	if u != nil {
		return u.UnmarshalCTYValue(value)
	}
	if uj != nil {
		bs, err := ctyValueToJSON(value)
		if err != nil {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type(), Detail: err}
		}
		return uj.UnmarshalJSON(bs)
	}
	if ut != nil {
		return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type()}
	}
	rv = pv
	rt := rv.Type()
	if value.IsNull() {
		rv.Set(reflect.Zero(rt))
		return nil
	}
	if rv.Kind() == reflect.Interface && rv.NumMethod() == 0 {
		converted, err := convertCTYObject(value)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(converted))
		return nil
	}
	if rv.Kind() == reflect.Map {
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rt))
		}
		if rt.Key().Kind() != reflect.String {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: rt}
		}
		valueMap := value.AsValueMap()
		for k, v := range valueMap {
			elemRv := reflect.New(rt.Elem())
			if err := unmarshalCTYValue(v, elemRv.Elem()); err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), elemRv.Elem())
		}
		return nil
	}
	if rv.Kind() == reflect.Struct {
		fields := getStructFileds(rt)
		valueMap := value.AsValueMap()
		for _, field := range fields {
			v, ok := valueMap[field.tagName]
			if !ok {
				continue
			}
			fv := rv
			for _, i := range field.index {
				fv = fv.Field(i)
			}
			if err := unmarshalCTYValue(v, fv); err != nil {
				return err
			}
		}
		return nil
	}
	return &UnmarshalTypeError{CTYType: value.Type(), Type: rt}
}

func unmarshalCTYPrimitive(value cty.Value, rv reflect.Value) error {
	if !value.IsKnown() {
		return &UnknownValueError{Value: value}
	}
	u, uj, ut, pv := indirect(rv, value.IsNull())
	if u != nil {
		return u.UnmarshalCTYValue(value)
	}
	if uj != nil {
		bs, err := ctyValueToJSON(value)
		if err != nil {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type(), Detail: err}
		}
		return uj.UnmarshalJSON(bs)
	}
	if ut != nil {
		if value.Type() != cty.String {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: rv.Type()}
		}
		if value.IsNull() {
			return nil
		}
		return ut.UnmarshalText([]byte(value.AsString()))
	}
	switch pv.Kind() {
	case reflect.Bool:
		if value.Type() != cty.Bool {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
		}
		pv.SetBool(value.True())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Type() != cty.Number {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
		}
		num, _ := value.AsBigFloat().Int64()
		pv.SetInt(num)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if value.Type() != cty.Number {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
		}
		num, _ := value.AsBigFloat().Uint64()
		pv.SetUint(num)
	case reflect.Float32, reflect.Float64:
		if value.Type() != cty.Number {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
		}
		num, _ := value.AsBigFloat().Float64()
		pv.SetFloat(num)
	case reflect.String:
		if value.Type() != cty.String {
			return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
		}
		pv.SetString(value.AsString())
	case reflect.Interface:
		if pv.NumMethod() == 0 {
			converted, err := convertCTYValue(value)
			if err != nil {
				return err
			}
			pv.Set(reflect.ValueOf(converted))
			return nil
		}
	default:
		return &UnmarshalTypeError{CTYType: value.Type(), Type: pv.Type()}
	}
	return nil
}

func ConvertCTYValue(value cty.Value) (any, error) {
	if !value.IsKnown() {
		return nil, &UnknownValueError{Value: value}
	}
	return convertCTYValue(value)
}

func convertCTYValue(value cty.Value) (any, error) {
	t := value.Type()
	switch {
	case t.IsListType() || t.IsTupleType() || t.IsSetType():
		return convertCTYList(value)
	case t.IsMapType() || t.IsObjectType():
		return convertCTYObject(value)
	case t.IsPrimitiveType():
		return convertCTYPrimitive(value)
	default:
		return nil, fmt.Errorf("hclutil: cannot convert %s to Go value", t.FriendlyName())
	}
}

func convertCTYList(value cty.Value) (any, error) {
	valueSlice := value.AsValueSlice()
	result := make([]any, len(valueSlice))
	for i, v := range valueSlice {
		var err error
		result[i], err = convertCTYValue(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func convertCTYObject(value cty.Value) (any, error) {
	valueMap := value.AsValueMap()
	result := make(map[string]any, len(valueMap))
	for k, v := range valueMap {
		var err error
		result[k], err = convertCTYValue(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func convertCTYPrimitive(value cty.Value) (any, error) {
	if value.IsNull() {
		return nil, nil
	}
	switch value.Type() {
	case cty.Bool:
		return value.True(), nil
	case cty.Number:
		return value.AsBigFloat(), nil
	case cty.String:
		return value.AsString(), nil
	default:
		return nil, fmt.Errorf("hclutil: cannot convert %s to Go value", value.Type().FriendlyName())
	}
}

// CTYValueUnmarshaler is an interface for types that can be decoded from a cty.Value.
type CTYValueUnmarshaler interface {
	UnmarshalCTYValue(cty.Value) error
}

type UnknownValueError struct {
	Value cty.Value
}

func (e *UnknownValueError) Error() string {
	return "hclutil: unknown value " + e.Value.GoString()
}

// InvalidUnmarshalError describes an invalid argument passed to UnmarshalCTYValue.
type InvalidUnmarshalError struct {
	Type reflect.Type
}

// Error implements the error interface.
func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "hclutil: UnmarshalCTYValue(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "hclutil: UnmarshalCTYValue(non-pointer " + e.Type.String() + ")"
	}
	return "hclutil: UnmarshalCTYValue(nil " + e.Type.String() + ")"
}

// UnmarshalTypeError describes a type missmatch between the cty.Type and the target type.
type UnmarshalTypeError struct {
	CTYType cty.Type
	Type    reflect.Type
	Detail  error
}

// Error implements the error interface.
func (e *UnmarshalTypeError) Error() string {
	return "hclutl: cannot unmarshal " + e.CTYType.FriendlyName() + " into Go value of type " + e.Type.String()
}

func (e *UnmarshalTypeError) Unwrap() error {
	return e.Detail
}
