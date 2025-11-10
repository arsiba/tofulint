// This file defines helper functions for encoding and decoding Terraform// variable files (tfvars) and HCL expressions.
package terraform

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// EncodeTfvarsFunc converts a cty value into a Terraform .tfvars formatted string.
var EncodeTfvarsFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowNull:        true,
			AllowDynamicType: true,
			AllowUnknown:     true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, _ cty.Type) (cty.Value, error) {
		// Validate argument count and type.
		if len(args) != 1 {
			return cty.NilVal, fmt.Errorf("exactly one argument is required")
		}
		v := args[0]
		t := v.Type()

		if v.IsNull() {
			return cty.NilVal, function.NewArgErrorf(1, "cannot encode null value as tfvars")
		}
		if !v.IsWhollyKnown() {
			return cty.UnknownVal(cty.String).RefineNotNull(), nil
		}

		// Collect map or object keys to ensure deterministic ordering.
		var keys []string
		switch {
		case t.IsObjectType():
			for k := range t.AttributeTypes() {
				keys = append(keys, k)
			}
		case t.IsMapType():
			for it := v.ElementIterator(); it.Next(); {
				k, _ := it.Element()
				keys = append(keys, k.AsString())
			}
		default:
			return cty.NilVal, function.NewArgErrorf(1, "expected an object or map for tfvars encoding")
		}
		sort.Strings(keys)

		f := hclwrite.NewEmptyFile()
		body := f.Body()

		for _, k := range keys {
			if !hclsyntax.ValidIdentifier(k) {
				return cty.NilVal, function.NewArgErrorf(1, "invalid variable name %q: must be a valid identifier", k)
			}
			val, _ := hcl.Index(v, cty.StringVal(k), nil)
			body.SetAttributeValue(k, val)
		}

		return cty.StringVal(string(f.Bytes())), nil
	},
})

// DecodeTfvarsFunc parses a tfvars-formatted string and returns its cty object representation.
var DecodeTfvarsFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:      "src",
			Type:      cty.String,
			AllowNull: true,
		},
	},
	Type: function.StaticReturnType(cty.DynamicPseudoType),
	Impl: func(args []cty.Value, _ cty.Type) (cty.Value, error) {
		// Validate argument count and type.
		if len(args) != 1 {
			return cty.NilVal, fmt.Errorf("exactly one argument is required")
		}
		arg := args[0]
		if arg.Type() != cty.String {
			return cty.NilVal, fmt.Errorf("argument must be a string")
		}
		if arg.IsNull() {
			return cty.NilVal, fmt.Errorf("cannot decode tfvars from a null value")
		}
		if !arg.IsKnown() {
			return cty.DynamicVal, nil
		}

		src := []byte(arg.AsString())
		file, diags := hclsyntax.ParseConfig(src, "<tfvars>", hcl.InitialPos)
		if diags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid tfvars syntax: %s", diags.Error())
		}

		attrs, diags := file.Body.JustAttributes()
		if diags.HasErrors() {
			return cty.NilVal, fmt.Errorf("invalid tfvars content: %s", diags.Error())
		}

		out := make(map[string]cty.Value, len(attrs))
		for name, attr := range attrs {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return cty.NilVal, fmt.Errorf("invalid expression for variable %q: %s", name, diags.Error())
			}
			out[name] = val
		}
		return cty.ObjectVal(out), nil
	},
})

// EncodeExprFunc converts a cty value into its canonical HCL expression string.
var EncodeExprFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "value",
			Type:             cty.DynamicPseudoType,
			AllowNull:        true,
			AllowDynamicType: true,
			AllowUnknown:     true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, _ cty.Type) (cty.Value, error) {
		// Validate arguments.
		if len(args) != 1 {
			return cty.NilVal, fmt.Errorf("exactly one argument is required")
		}
		v := args[0]

		if !v.IsWhollyKnown() {
			ret := cty.UnknownVal(cty.String).RefineNotNull()

			if !v.Range().CouldBeNull() {
				switch ty := v.Type(); {
				case ty.IsObjectType() || ty.IsMapType():
					ret = ret.Refine().StringPrefixFull("{").NewValue()
				case ty.IsTupleType() || ty.IsListType() || ty.IsSetType():
					ret = ret.Refine().StringPrefixFull("[").NewValue()
				case ty == cty.String:
					ret = ret.Refine().StringPrefixFull(`"`).NewValue()
				}
			}
			return ret, nil
		}

		src := bytes.TrimSpace(hclwrite.TokensForValue(v).Bytes())
		return cty.StringVal(string(src)), nil
	},
})
