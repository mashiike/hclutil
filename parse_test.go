package hclutil_test

import (
	"bytes"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/mashiike/hclutil"
	"github.com/stretchr/testify/require"
)

func TestParse__HCLFile(t *testing.T) {
	t.Parallel()
	body, writer, diags := hclutil.Parse("testdata/hcl_file.hcl")
	diagsReport(t, diags)
	require.NotNil(t, body)
	var buf bytes.Buffer
	writer.SetOutput(&buf)
	err := writer.WriteDiagnostics(hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "test error",
			Detail:   "test error detail",
			Subject:  body.MissingItemRange().Ptr(),
		},
	})
	require.EqualError(t, err, "diagnostics had errors, see above for details")
	require.Equal(t, "Error: test error\n\n  on testdata/hcl_file.hcl line 1:\n   1: text = \"hoge\"\n\ntest error detail\n\n", buf.String())
	attrs, _ := body.JustAttributes()
	attrKeys := make([]string, 0, len(attrs))
	for k := range attrs {
		attrKeys = append(attrKeys, k)
	}
	require.ElementsMatch(t, []string{"boolean", "text", "number"}, attrKeys)
}

func TestParse__HCLFile__NotFound(t *testing.T) {
	t.Parallel()
	_, _, diags := hclutil.Parse("testdata/notfound.hcl")
	require.EqualError(t, diags, "<nil>: Parse failed; stat testdata/notfound.hcl: no such file or directory")
}

func TestParse__JSONFile(t *testing.T) {
	t.Parallel()
	body, writer, diags := hclutil.Parse("testdata/json_file.hcl.json")
	diagsReport(t, diags)
	require.NotNil(t, body)
	var buf bytes.Buffer
	writer.SetOutput(&buf)
	err := writer.WriteDiagnostics(hcl.Diagnostics{
		{
			Severity: hcl.DiagError,
			Summary:  "test error",
			Detail:   "test error detail",
			Subject:  body.MissingItemRange().Ptr(),
		},
	})
	require.EqualError(t, err, "diagnostics had errors, see above for details")
	require.Equal(t, "Error: test error\n\n  on testdata/json_file.hcl.json line 5:\n   5: }\n\ntest error detail\n\n", buf.String())
	attrs, _ := body.JustAttributes()
	attrKeys := make([]string, 0, len(attrs))
	for k := range attrs {
		attrKeys = append(attrKeys, k)
	}
	require.ElementsMatch(t, []string{"boolean", "text", "number"}, attrKeys)
}

func TestParse__JSONFile__NotFound(t *testing.T) {
	t.Parallel()
	_, _, diags := hclutil.Parse("testdata/notfound.hcl.json")
	require.EqualError(t, diags, "<nil>: Parse failed; stat testdata/notfound.hcl.json: no such file or directory")
}

func TestParse__Dir(t *testing.T) {
	t.Parallel()
	body, writer, diags := hclutil.Parse("testdata/simple")
	diagsReport(t, diags)
	require.NotNil(t, body)
	require.ElementsMatch(t, []string{
		"testdata/simple/hcl_file.hcl",
		"testdata/simple/json_file.hcl.json",
	}, writer.Files())
	attrs, _ := body.JustAttributes()
	attrKeys := make([]string, 0, len(attrs))
	for k := range attrs {
		attrKeys = append(attrKeys, k)
	}
	require.ElementsMatch(t, []string{"boolean", "text", "number"}, attrKeys)
}

func TestParse__Dir__NotFound(t *testing.T) {
	t.Parallel()
	_, _, diags := hclutil.Parse("testdata/notfound")
	require.EqualError(t, diags, "<nil>: Parse failed; stat testdata/notfound: no such file or directory")
}
