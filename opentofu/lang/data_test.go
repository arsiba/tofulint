// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lang

import (
	"context"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/arsiba/tofulint/opentofu/addrs"
)

type dataForTests struct {
	CountAttrs              map[string]cty.Value
	ForEachAttrs            map[string]cty.Value
	LocalValues             map[string]cty.Value
	PathAttrs               map[string]cty.Value
	TerraformAttrs          map[string]cty.Value
	InputVariables          map[string]cty.Value
	StaticValidateReference map[string]cty.Value
}

func (d *dataForTests) StaticValidateReferences(ctx context.Context, refs []*addrs.Reference, self addrs.Referenceable, source addrs.Referenceable) hcl.Diagnostics {
	return nil
}

var _ Data = &dataForTests{}

func (d *dataForTests) GetCountAttr(ctx context.Context, addr addrs.CountAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.CountAttrs[addr.Name], nil
}

func (d *dataForTests) GetForEachAttr(ctx context.Context, addr addrs.ForEachAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.ForEachAttrs[addr.Name], nil
}

func (d *dataForTests) GetInputVariable(ctx context.Context, addr addrs.InputVariable, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.InputVariables[addr.Name], nil
}

func (d *dataForTests) GetLocalValue(ctx context.Context, addr addrs.LocalValue, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.LocalValues[addr.Name], nil
}

func (d *dataForTests) GetPathAttr(ctx context.Context, addr addrs.PathAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.PathAttrs[addr.Name], nil
}

func (d *dataForTests) GetTerraformAttr(ctx context.Context, addr addrs.TerraformAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return d.TerraformAttrs[addr.Name], nil
}
