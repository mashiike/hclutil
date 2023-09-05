# hclutil
HCL utility for Golang


[![GoDoc](https://godoc.org/github.com/mashiike/hclutil?status.svg)](https://godoc.org/github.com/mashiike/hclutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/mashiike/hclutil)](https://goreportcard.com/report/github.com/mashiike/hclutil)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Overview

This package provides middleware and utility functions for easy logging management. It enables color-coded display for different log levels and automatically collects attributes set in the context. This allows developers to have flexible logging recording and analysis capabilities.

## Installation

```bash
go get github.com/mashiike/hclutil
```

## Usage

sample code

```go
package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	App struct {
		RequiredVresion cty.Value `hcl:"required_version"`
		Name            string    `hcl:"name"`
		Description     string    `hcl:"description"`
	} `hcl:"app,block"`
}

func main() {
	var v Config
	body, writer, diags := hclutil.Parse("./")
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}
	evalCtx := hclutil.NewEvalContext()
	body, evalCtx, diags = hclutil.DecodeLocals(body, evalCtx)
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}
	diags = gohcl.DecodeBody(body, evalCtx, &v)
	if diags.HasErrors() {
		writer.WriteDiagnostics(diags)
		log.Fatal("parse failed")
	}
	fmt.Println("name:", v.App.Name)
	fmt.Println("description:", v.App.Description)

	var vc hclutil.VersionConstraints
	if err := hclutil.UnmarshalCTYValue(v.App.RequiredVresion, &vc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("required_version:", vc.String())
	fmt.Println("v1.2.3 is satisfied:", vc.ValidateVersion("v1.2.3") == nil)
}
```

### NewEvalContext

this function is create new EvalContext with helpful functions.

### DecodeLocals

this function is decode locals block and return new body and EvalContext.

### UnmarshalCTYValue

this function is unmarshal cty.Value to Any.

## License
This project is licensed under the MIT License - see the LICENSE(./LICENCE) file for details.

## Contribution
Contributions, bug reports, and feature requests are welcome. Pull requests are also highly appreciated. For more details, please
