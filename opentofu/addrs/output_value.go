// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package addrs

// OutputValue is the address of an output value, in the context of the module
// that is defining it.
//
// This is related to but separate from ModuleCallOutput, which represents
// a module output from the perspective of its parent module. Outputs are
// referenceable from the testing scope, in general tofu operation users
// will be referencing ModuleCallOutput.
type OutputValue struct {
	referenceable
	Name string
}

func (v OutputValue) String() string {
	return "output." + v.Name
}
