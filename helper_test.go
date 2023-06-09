package hclutil_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
)

func diagsReport(t *testing.T, diags hcl.Diagnostics) {
	t.Helper()
	for _, diag := range diags {
		if diag.Severity == hcl.DiagError {
			t.Error("ERROR:", diag)
		}
		if diag.Severity == hcl.DiagWarning {
			t.Log("WARN:", diag)
		}
	}
}
