package hclutil

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// ExtractAttributes interprets the contents of the HCL body as attributes,
// allowing for the contents to be accessed without prior knowledge of the structure.
// This function ignores any blocks present in the body and focuses solely on extracting attributes.
func ExtructAttributes(body hcl.Body) (hcl.Attributes, hcl.Diagnostics) {
	attrs, d := body.JustAttributes()
	diags := make(hcl.Diagnostics, 0, len(d))
	for _, diag := range d {
		if strings.Contains(diag.Detail, "Blocks are not allowed here") {
			continue
		}
		diags = diags.Append(diag)
	}
	return attrs, diags
}
