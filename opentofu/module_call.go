// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package opentofu

import (
	"fmt"

	"github.com/arsiba/tofulint/opentofu/addrs"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// ModuleCall represents a "module" block in a module or file.
type ModuleCall struct {
	Name string

	SourceAddr    addrs.ModuleSource
	SourceAddrRaw string

	DeclRange hcl.Range
}

func decodeModuleBlock(block *hclext.Block) (*ModuleCall, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	mc := &ModuleCall{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}

	if attr, exists := block.Body.Attributes["source"]; exists {
		valDiags := gohcl.DecodeExpression(attr.Expr, nil, &mc.SourceAddrRaw)
		diags = append(diags, valDiags...)

		if !diags.HasErrors() {
			var err error
			mc.SourceAddr, err = addrs.ParseModuleSource(mc.SourceAddrRaw)
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Invalid module source address",
					Detail:   fmt.Sprintf("Failed to parse module source address: %s", err),
					Subject:  attr.Expr.Range().Ptr(),
				})
			}
		}
	}

	return mc, diags
}

var moduleBlockSchema = &hclext.BodySchema{
	Attributes: []hclext.AttributeSchema{
		{
			Name: "source",
		},
	},
}

// CallModuleType is a type of module to call.
// This is primarily used to control module walker behavior.
type CallModuleType int32

const (
	// CallAllModule calls all (local/remote) modules.
	CallAllModule CallModuleType = iota

	// CallLocalModule calls only local modules.
	CallLocalModule

	// CallNoModule does not call any modules.
	CallNoModule
)

func AsCallModuleType(s string) (CallModuleType, error) {
	switch s {
	case "all":
		return CallAllModule, nil
	case "local":
		return CallLocalModule, nil
	case "none":
		return CallNoModule, nil
	default:
		return CallAllModule, fmt.Errorf("%s is invalid call module type. Allowed values are: all, local, none", s)
	}
}

func (c CallModuleType) String() string {
	switch c {
	case CallAllModule:
		return "all"
	case CallLocalModule:
		return "local"
	case CallNoModule:
		return "none"
	default:
		panic("never happened")
	}
}
