package hclutil_test

import (
	"errors"
	"testing"

	"github.com/mashiike/hclutil"
	"github.com/zclconf/go-cty/cty"
)

func TestVersionConstraints(t *testing.T) {
	t.Parallel()
	var vc hclutil.VersionConstraints
	t.Run("UnmarshalCTYValue", func(t *testing.T) {
		err := hclutil.UnmarshalCTYValue(cty.StringVal(">= 1.2.3"), &vc)
		if err != nil {
			t.Error(err)
		}
		if vc.String() != ">= 1.2.3" {
			t.Errorf("vc.String() = %s, want >= 1.2.3", vc.String())
		}
	})
	t.Run("MarshalCTYValue", func(t *testing.T) {
		v, err := hclutil.MarshalCTYValue(&vc)
		if err != nil {
			t.Error(err)
		}
		if v.Type() != cty.String {
			t.Errorf("v.Type() = %s, want string", v.Type())
			t.FailNow()
		}
		if v.AsString() != ">= 1.2.3" {
			t.Errorf("v = %s, want >= 1.2.3", v.AsString())
		}
	})
}

func TestVersionConstraints__ValidateVersion(t *testing.T) {
	t.Parallel()
	var vc hclutil.VersionConstraints
	err := hclutil.UnmarshalCTYValue(cty.StringVal(">= 1.2.3"), &vc)
	if err != nil {
		t.Error(err)
	}
	t.Run("satisfied", func(t *testing.T) {
		err := vc.ValidateVersion("1.2.3")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("not satisfied", func(t *testing.T) {
		err := vc.ValidateVersion("1.2.2")
		if err == nil {
			t.Error("err == nil, want error")
		}
		var ve *hclutil.VersionConstraintNotSatisfiedError
		if !errors.As(err, &ve) {
			t.Errorf("err = %T, want VersionConstraintNotSatisfiedError", err)
		}
	})
	t.Run("invalid version", func(t *testing.T) {
		err := vc.ValidateVersion("current")
		if err == nil {
			t.Error("err == nil, want error")
		}
		var ve *hclutil.InvalidVersionError
		if !errors.As(err, &ve) {
			t.Errorf("err = %T, want InvalidVersionError", err)
		}
	})
	t.Run("no constraint", func(t *testing.T) {
		var vc hclutil.VersionConstraints
		err := vc.ValidateVersion("1.2.3")
		if err != nil {
			t.Error(err)
		}
	})
}
