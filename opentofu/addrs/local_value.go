// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package addrs

// LocalValue is the address of a local value.
type LocalValue struct {
	referenceable
	Name string
}

func (v LocalValue) String() string {
	return "local." + v.Name
}
