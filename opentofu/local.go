package opentofu

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// A consistent detail message for all "not a valid identifier" diagnostics.
const badIdentifierDetail = "A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes."

// Local represents a single entry from a "locals" block in a module or file.
// The "locals" block itself is not represented, because it serves only to
// provide context for us to interpret its contents.
type Local struct {
	Name string
	Expr hcl.Expression

	DeclRange hcl.Range
}

func decodeLocalsBlock(block *hclext.Block) []*Local {

	locals := make([]*Local, 0, len(block.Body.Attributes))
	for name, attr := range block.Body.Attributes {
		locals = append(locals, &Local{
			Name:      name,
			Expr:      attr.Expr,
			DeclRange: attr.Range,
		})
	}
	return locals
}

var localBlockSchema = &hclext.BodySchema{
	Mode: hclext.SchemaJustAttributesMode,
}
