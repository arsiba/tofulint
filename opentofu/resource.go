// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package opentofu

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// Resource represents a "resource" or "data" block in a module or file.
type Resource struct {
	Name string
	Type string

	DeclRange hcl.Range
	TypeRange hcl.Range
}

func decodeResourceBlock(block *hclext.Block) *Resource {
	r := &Resource{
		Type:      block.Labels[0],
		Name:      block.Labels[1],
		DeclRange: block.DefRange,
		TypeRange: block.LabelRanges[0],
	}

	return r
}
