package hclutil

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// BlockRestrictionSchema is a schema for block restriction
type BlockRestrictionSchema struct {
	// Type is the block type
	Type string
	// Required is the block required
	Required bool
	// Unique is the block unique
	Unique bool
	// UniqueLabels is the block unique labels
	UniqueLabels bool
}

// RestrictBlock implements the restriction that block
func RestrictBlock(content *hcl.BodyContent, schemas ...BlockRestrictionSchema) hcl.Diagnostics {
	var diags hcl.Diagnostics
	blocksByType := content.Blocks.ByType()
	for _, schema := range schemas {
		blocks, ok := blocksByType[schema.Type]
		if !ok {
			if schema.Required {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf(`Missing "%s" block`, schema.Type),
					Detail:   fmt.Sprintf(`A "%s" block is required.`, schema.Type),
					Subject:  content.MissingItemRange.Ptr(),
				})
			}
			continue
		}
		if schema.Unique {
			if len(blocks) > 1 {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  fmt.Sprintf(`Duplicate "%s" block`, schema.Type),
					Detail:   fmt.Sprintf(`Only one "%s" block is allowed. Another was defined at %s`, schema.Type, blocks[1].DefRange.String()),
					Subject:  blocks[0].DefRange.Ptr(),
				})
			}
		}
		if schema.UniqueLabels {
			fqns := make(map[string]*hcl.Range, len(blocks))
			for _, block := range blocks {
				fqn := block.Type + "." + strings.Join(block.Labels, ".")
				if r, ok := fqns[fqn]; ok {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf(`Duplicate "%s" block`, fqn),
						Detail:   fmt.Sprintf(`Only one "%s" block is allowed. Another was defined at %s`, fqn, r.String()),
						Subject:  block.DefRange.Ptr(),
					})
				} else {
					fqns[fqn] = block.DefRange.Ptr()
				}
			}
		}
	}
	return diags
}
