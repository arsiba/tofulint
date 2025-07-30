// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lang

import (
	"sync"

	"github.com/zclconf/go-cty/cty/function"

	"github.com/arsiba/tofulint/opentofu/addrs"
)

// Scope is the main type in this package, allowing dynamic evaluation of
// blocks and expressions based on some contextual information that informs
// which variables and functions will be available.
type Scope struct {
	// Data is used to resolve references in expressions.
	Data Data

	// SelfAddr is the address that the "self" object should be an alias of,
	// or nil if the "self" object should not be available at all.
	SelfAddr addrs.Referenceable

	// SourceAddr is the address of the source item for the scope. This will
	// affect any scoped resources that can be accessed from within this scope.
	//
	// If nil, access is assumed to be at the module level. So, in practice this
	// only needs to be set for items that should be able to access something
	// hidden in their own scope.
	SourceAddr addrs.Referenceable

	// BaseDir is the base directory used by any interpolation functions that
	// accept filesystem paths as arguments.
	BaseDir string

	// PureOnly can be set to true to request that any non-pure functions
	// produce unknown value results rather than actually executing. This is
	// important during a plan phase to avoid generating results that could
	// then differ during apply.
	PureOnly bool

	funcsLock sync.Mutex
	funcs     map[string]function.Function
}
