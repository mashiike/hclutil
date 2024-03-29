package hclutil

import (
	"reflect"
	"strings"
	"sync"
)

var fieldCache sync.Map

type field struct {
	name      string
	typ       reflect.Type
	hclTag    string
	ctyTag    string
	tagName   string
	index     []int
	omitEmpty bool
}

type structFields []field

func getStructFileds(rt reflect.Type) structFields {
	if f, ok := fieldCache.Load(rt); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(rt, newStructFields(rt))
	return f.(structFields)
}

func newStructFields(rt reflect.Type) structFields {
	var fields structFields
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		ctyTag := f.Tag.Get("cty")
		hclTag := f.Tag.Get("hcl")
		name := strings.Split(ctyTag, ",")[0]
		if name == "-" {
			continue
		}
		if name == "" {
			name = strings.Split(hclTag, ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = camelcaseToSnakecase(f.Name)
			}
		}
		omitEmpty := strings.Contains(ctyTag, ",omitempty")
		if ctyTag == "" {
			omitEmpty = strings.Contains(hclTag, ",omitempty")
		}
		if f.Anonymous && ft.Kind() == reflect.Struct {
			embeddedFields := getStructFileds(ft)
			for _, embeddedField := range embeddedFields {
				embeddedField.index = append(f.Index, embeddedField.index...)
				fields = append(fields, embeddedField)
			}
			continue
		}
		if !f.Anonymous && name != "" && f.IsExported() {
			fields = append(fields, field{
				name:      f.Name,
				typ:       f.Type,
				hclTag:    hclTag,
				ctyTag:    ctyTag,
				tagName:   name,
				index:     []int{i},
				omitEmpty: omitEmpty,
			})
		}
	}
	return fields
}

func camelcaseToSnakecase(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}
