package hclutil

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

func TraversalToString(t hcl.Traversal) string {
	parts := make([]string, 0, len(t))
	if t.IsRelative() {
		parts = append(parts, "")
	}
	for _, tr := range t {
		parts = traverserToString(parts, tr)
	}

	return strings.Join(parts, ".")
}

func traverserToString(parts []string, tr hcl.Traverser) []string {
	switch tr := tr.(type) {
	case hcl.TraverseAttr:
		parts = append(parts, tr.Name)
	case hcl.TraverseIndex:
		var indexJson json.RawMessage
		indexStr := "?"
		if err := UnmarshalCTYValue(tr.Key, &indexJson); err == nil {
			indexStr = string(indexJson)
		}
		if len(parts) > 0 && parts[len(parts)-1] != "" {
			parts[len(parts)-1] += "[" + indexStr + "]"
		} else {
			parts = append(parts, "["+indexStr+"]")
		}
	case hcl.TraverseRoot:
		parts = append(parts, tr.Name)
	case hcl.TraverseSplat:
		for _, tt := range tr.Each {
			parts = traverserToString(parts, tt)
		}
	}
	return parts
}

// VariablesReffarances returns variables reffarance string list in expression
func VariablesReffarances(expr hcl.Expression) []string {
	travarsals := expr.Variables()
	vars := make([]string, len(travarsals))
	for i, t := range travarsals {
		vars[i] = TraversalToString(t)
	}
	return vars
}
