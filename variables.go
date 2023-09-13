package hclutil

import "github.com/zclconf/go-cty/cty"

// MergeVariables merges multiple variables into one.
func MergeVariables(vars ...map[string]cty.Value) map[string]cty.Value {
	var result map[string]cty.Value
	for _, v := range vars {
		result = mergeVariables(result, v)
	}
	return result
}

func mergeVariables(dst map[string]cty.Value, src map[string]cty.Value) map[string]cty.Value {
	if dst == nil {
		dst = make(map[string]cty.Value, len(src))
	}
	for key, value := range src {
		dstValue, ok := dst[key]
		if !ok {
			dst[key] = value
			continue
		}
		if !dstValue.Type().IsObjectType() {
			dst[key] = value
			continue
		}
		if !value.Type().IsObjectType() {
			dst[key] = value
			continue
		}
		dstValueMap := dstValue.AsValueMap()
		srcValueMap := value.AsValueMap()
		dst[key] = cty.ObjectVal(mergeVariables(dstValueMap, srcValueMap))
	}
	return dst
}
