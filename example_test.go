package hclutil_test

import (
	"fmt"
	"log"
	"testing/fstest"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func Example() {
	var sysValue int
	config := `

locals {
	value = "hoge"
}

app {
	name = upper(local.value)
	description = templatefile("description.txt", { value = local.value })
}
	`
	testFs := fstest.MapFS{
		"description.txt": {
			Data:    []byte("hello, ${value}"),
			Mode:    0,
			ModTime: time.Now(),
			Sys:     &sysValue,
		},
		"config.hcl": {
			Data:    []byte(config),
			Mode:    0,
			ModTime: time.Now(),
			Sys:     &sysValue,
		},
	}

	body, writer, diags := hclutil.ParseFS(testFs)
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}
	evalCtx := hclutil.NewEvalContext(hclutil.WithFS(testFs))
	body, evalCtx, diags = hclutil.DecodeLocals(body, evalCtx)
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}

	var v struct {
		App struct {
			Name        string `hcl:"name"`
			Description string `hcl:"description"`
		} `hcl:"app,block"`
	}
	diags = gohcl.DecodeBody(body, evalCtx, &v)
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}
	fmt.Println("name:", v.App.Name)
	fmt.Println("description:", v.App.Description)
	// Output:
	// name: HOGE
	// description: hello, hoge
}

func ExampleMarshalCTYValue() {
	var v struct {
		Name string `cty:"name"`
		Age  int    `cty:"age"`
	}
	v.Name = "Alice"
	v.Age = 12

	ctyValue, err := hclutil.MarshalCTYValue(v)
	if err != nil {
		log.Fatalf("failed to marshal cty value: %s", err)
	}
	fmt.Println(hclutil.MustDumpCtyValue(ctyValue))

	// Output:
	// {"age":12,"name":"Alice"}
}

func ExampleUnmarshalCTYValue() {
	var v struct {
		Name string `cty:"name"`
		Age  int    `cty:"age"`
	}

	ctyValue := cty.ObjectVal(map[string]cty.Value{
		"name": cty.StringVal("Alice"),
		"age":  cty.NumberIntVal(12),
	})

	if err := hclutil.UnmarshalCTYValue(ctyValue, &v); err != nil {
		log.Fatalf("failed to unmarshal cty value: %s", err)
	}
	fmt.Println("name:", v.Name)
	fmt.Println("age:", v.Age)

	// Output:
	// name: Alice
	// age: 12
}

func ExampleVersionConstraints() {
	var vc hclutil.VersionConstraints
	expr, diags := hclsyntax.ParseExpression([]byte(`">= 1.0.0, < 2.0.0"`), "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		log.Fatal("parse failed")
	}
	diags = vc.DecodeExpression(expr, hclutil.NewEvalContext())
	if diags.HasErrors() {
		log.Fatal("parse failed")
	}
	fmt.Println(vc)
	fmt.Println("v1.2.0 is satisfied:", vc.IsSutisfied("v1.2.0"))
	// Output:
	// {>= 1.0.0, < 2.0.0}
	// v1.2.0 is satisfied: true
}
