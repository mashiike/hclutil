package hclutil

import (
	"errors"
	"fmt"

	goVersion "github.com/hashicorp/go-version"
	"github.com/zclconf/go-cty/cty"
)

type VersionConstraints struct {
	goVersion.Constraints
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