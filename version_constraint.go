package hclutil

import (
	"errors"
	"fmt"

	goVersion "github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// VersionConstraints is a wrapper of goVersion.Constraints to implement CTYValueMarshaler and CTYValueUnmarshaler interface.
type VersionConstraints struct {
	goVersion.Constraints
}

func (vc *VersionConstraints) String() string {
	if vc == nil {
		return ""
	}
	return vc.Constraints.String()
}

func (vc *VersionConstraints) DecodeExpression(expr hcl.Expression, ctx *hcl.EvalContext) hcl.Diagnostics {
	v, diags := expr.Value(ctx)
	if diags.HasErrors() {
		return diags
	}
	if err := vc.UnmarshalCTYValue(v); err != nil {
		return hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Invalid version constraint",
			Detail:   err.Error(),
		}}
	}
	return nil
}

func (vc *VersionConstraints) UnmarshalCTYValue(value cty.Value) error {
	if !value.IsKnown() {
		return errors.New("version constraint must be known")
	}
	if value.IsNull() {
		return nil
	}
	if value.Type() != cty.String {
		return fmt.Errorf("version constraint expected to be a string, got %s", value.Type().FriendlyName())
	}
	constraints, err := goVersion.NewConstraint(value.AsString())
	if err != nil {
		return err
	}
	vc.Constraints = constraints
	return nil
}

func (vc *VersionConstraints) MarshalCTYValue() (cty.Value, error) {
	if vc == nil {
		return cty.NullVal(cty.String), nil
	}
	return cty.StringVal(vc.String()), nil
}

// InvalidVersionError is an error type for invalid version.
type InvalidVersionError struct {
	Version string
	err     error
}

func (err *InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version %s: %s", err.Version, err.err)
}

func (err *InvalidVersionError) Unwrap() error {
	return err.err
}

// VersionConstraintNotSatisfiedError is an error type for version constraint not satisfied.
type VersionConstraintNotSatisfiedError struct {
	Constraint *VersionConstraints
	Version    string
}

func (err *VersionConstraintNotSatisfiedError) Error() string {
	return fmt.Sprintf("version %s does not satisfy constraint %s", err.Version, err.Constraint.String())
}

func (vc *VersionConstraints) ValidateVersion(str string) error {
	if vc == nil {
		return nil
	}
	v, err := goVersion.NewVersion(str)
	if err != nil {
		return &InvalidVersionError{Version: str, err: err}
	}
	if !vc.Check(v) {
		return &VersionConstraintNotSatisfiedError{Constraint: vc, Version: str}
	}
	return nil
}

func (vc *VersionConstraints) IsSutisfied(str string) bool {
	if vc == nil {
		return true
	}
	v, err := goVersion.NewVersion(str)
	if err != nil {
		return false
	}
	return vc.Check(v)
}
