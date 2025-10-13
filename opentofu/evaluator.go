package opentofu

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agext/levenshtein"
	"github.com/arsiba/tofulint/opentofu/addrs"
	"github.com/arsiba/tofulint/opentofu/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/arsiba/tofulint-plugin-sdk/hclext"
	"github.com/arsiba/tofulint-plugin-sdk/terraform/lang/marks"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type ContextMeta struct {
	Env                string
	OriginalWorkingDir string
}

type CallStack struct {
	addrs map[string]addrs.Reference
	stack []string
}

func NewCallStack() *CallStack {
	return &CallStack{
		addrs: make(map[string]addrs.Reference),
		stack: make([]string, 0),
	}
}

func (g *CallStack) Push(addr addrs.Reference) hcl.Diagnostics {
	g.stack = append(g.stack, addr.Subject.String())

	if _, exists := g.addrs[addr.Subject.String()]; exists {
		return hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "circular reference found",
				Detail:   g.String(),
				Subject:  addr.SourceRange.Ptr(),
			},
		}
	}
	g.addrs[addr.Subject.String()] = addr
	return hcl.Diagnostics{}
}

func (g *CallStack) Pop() {
	if g.Empty() {
		panic("cannot pop from empty stack")
	}

	addr := g.stack[len(g.stack)-1]
	g.stack = g.stack[:len(g.stack)-1]
	delete(g.addrs, addr)
}

func (g *CallStack) String() string {
	return strings.Join(g.stack, " -> ")
}

func (g *CallStack) Empty() bool {
	return len(g.stack) == 0
}

func (g *CallStack) Clear() {
	g.addrs = make(map[string]addrs.Reference)
	g.stack = make([]string, 0)
}

type Evaluator struct {
	Meta           *ContextMeta
	ModulePath     addrs.ModuleInstance
	Config         *Config
	VariableValues map[string]map[string]cty.Value
	CallStack      *CallStack
}

func (e *Evaluator) EvaluateExpr(expr hcl.Expression, wantType cty.Type) (cty.Value, hcl.Diagnostics) {
	if e == nil {
		panic("evaluator must not be nil")
	}
	return e.scope().EvalExpr(expr, wantType)
}

func (e *Evaluator) ExpandBlock(body hcl.Body, schema *hclext.BodySchema) (hcl.Body, hcl.Diagnostics) {
	if e == nil {
		return body, nil
	}
	return e.scope().ExpandBlock(body, schema)
}

type evaluationData struct {
	Evaluator  *Evaluator
	ModulePath addrs.ModuleInstance
}

func (d *evaluationData) StaticValidateReferences(ctx context.Context, refs []*addrs.Reference, self addrs.Referenceable, source addrs.Referenceable) hcl.Diagnostics {
	return nil
}

func (e *Evaluator) scope() *lang.Scope {
	return &lang.Scope{
		Data: &evaluationData{
			Evaluator:  e,
			ModulePath: e.ModulePath,
		},
	}
}

var _ lang.Data = (*evaluationData)(nil)

func (d *evaluationData) GetCountAttr(ctx context.Context, addr addrs.CountAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return cty.UnknownVal(cty.Number), nil
}

func (d *evaluationData) GetForEachAttr(ctx context.Context, addr addrs.ForEachAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	return cty.DynamicVal, nil
}

func (d *evaluationData) GetInputVariable(ctx context.Context, addr addrs.InputVariable, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	moduleConfig := d.Evaluator.Config.DescendentForInstance(d.ModulePath)
	if moduleConfig == nil {
		panic(fmt.Sprintf("input variable read from %s, which has no configuration", d.ModulePath))
	}

	config := moduleConfig.Module.Variables[addr.Name]
	if config == nil {
		var suggestions []string
		for k := range moduleConfig.Module.Variables {
			suggestions = append(suggestions, k)
		}
		suggestion := nameSuggestion(addr.Name, suggestions)
		if suggestion != "" {
			suggestion = fmt.Sprintf(" Did you mean %q?", suggestion)
		} else {
			suggestion = fmt.Sprintf(" This variable can be declared with a variable %q {} block.", addr.Name)
		}

		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Reference to undeclared input variable`,
			Detail:   fmt.Sprintf(`An input variable with the name %q has not been declared.%s`, addr.Name, suggestion),
			Subject:  rng.Ptr(),
		})
		return cty.DynamicVal, diags
	}

	moduleAddrStr := d.ModulePath.String()
	vals := d.Evaluator.VariableValues[moduleAddrStr]
	if vals == nil {
		return cty.UnknownVal(config.Type), diags
	}

	val, isSet := vals[addr.Name]
	switch {
	case !isSet:
		val = config.Default
	case val.IsNull() && !config.Nullable && config.Default != cty.NilVal:
		val = config.Default
	}

	if config.TypeDefaults != nil && !val.IsNull() {
		val = config.TypeDefaults.Apply(val)
	}

	var err error
	val, err = convert.Convert(val, config.ConstraintType)
	if err != nil {
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Incorrect variable type`,
			Detail:   fmt.Sprintf(`The resolved value of variable %q is not appropriate: %s.`, addr.Name, err),
			Subject:  &config.DeclRange,
		})
		val = cty.UnknownVal(config.Type)
	}

	if config.Sensitive {
		val = val.Mark(marks.Sensitive)
	}

	return val, diags
}

func (d *evaluationData) GetLocalValue(ctx context.Context, addr addrs.LocalValue, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics

	moduleConfig := d.Evaluator.Config.DescendentForInstance(d.ModulePath)
	if moduleConfig == nil {
		panic(fmt.Sprintf("local value read from %s, which has no configuration", d.ModulePath))
	}

	config := moduleConfig.Module.Locals[addr.Name]
	if config == nil {
		var suggestions []string
		for k := range moduleConfig.Module.Locals {
			suggestions = append(suggestions, k)
		}
		suggestion := nameSuggestion(addr.Name, suggestions)
		if suggestion != "" {
			suggestion = fmt.Sprintf(" Did you mean %q?", suggestion)
		}

		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Reference to undeclared local value`,
			Detail:   fmt.Sprintf(`A local value with the name %q has not been declared.%s`, addr.Name, suggestion),
			Subject:  rng.Ptr(),
		})
		return cty.DynamicVal, diags
	}

	if diags := d.Evaluator.CallStack.Push(addrs.Reference{Subject: addr, SourceRange: rng}); diags.HasErrors() {
		return cty.UnknownVal(cty.DynamicPseudoType), diags
	}

	val, diags := d.Evaluator.EvaluateExpr(config.Expr, cty.DynamicPseudoType)
	d.Evaluator.CallStack.Pop()
	return val, diags
}

func (d *evaluationData) GetPathAttr(ctx context.Context, addr addrs.PathAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	switch addr.Name {
	case "cwd":
		var err error
		var wd string
		if d.Evaluator.Meta != nil {
			wd = d.Evaluator.Meta.OriginalWorkingDir
		}
		if wd == "" {
			wd, err = os.Getwd()
			if err != nil {
				diags = diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  `Failed to get working directory`,
					Detail:   fmt.Sprintf(`System error: %s`, err),
					Subject:  rng.Ptr(),
				})
				return cty.DynamicVal, diags
			}
		}
		wd, err = filepath.Abs(wd)
		if err != nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  `Failed to get working directory`,
				Detail:   fmt.Sprintf(`System error: %s`, err),
				Subject:  rng.Ptr(),
			})
			return cty.DynamicVal, diags
		}
		return cty.StringVal(filepath.ToSlash(wd)), diags

	case "module":
		moduleConfig := d.Evaluator.Config.DescendentForInstance(d.ModulePath)
		if moduleConfig == nil {
			panic(fmt.Sprintf("module.path read from module %s, which has no configuration", d.ModulePath))
		}
		sourceDir := moduleConfig.Module.SourceDir
		return cty.StringVal(filepath.ToSlash(sourceDir)), diags

	case "root":
		sourceDir := d.Evaluator.Config.Module.SourceDir
		return cty.StringVal(filepath.ToSlash(sourceDir)), diags

	default:
		suggestion := nameSuggestion(addr.Name, []string{"cwd", "module", "root"})
		if suggestion != "" {
			suggestion = fmt.Sprintf(" Did you mean %q?", suggestion)
		}
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Invalid "path" attribute`,
			Detail:   fmt.Sprintf(`The "path" object does not have an attribute named %q.%s`, addr.Name, suggestion),
			Subject:  rng.Ptr(),
		})
		return cty.DynamicVal, diags
	}
}

func (d *evaluationData) GetTerraformAttr(ctx context.Context, addr addrs.TerraformAttr, rng hcl.Range) (cty.Value, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	switch addr.Name {
	case "workspace":
		workspaceName := d.Evaluator.Meta.Env
		return cty.StringVal(workspaceName), diags

	case "env":
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Invalid "terraform" attribute`,
			Detail:   `The terraform.env attribute was deprecated and removed. Use terraform.workspace instead.`,
			Subject:  rng.Ptr(),
		})
		return cty.DynamicVal, diags

	default:
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  `Invalid "terraform" attribute`,
			Detail:   fmt.Sprintf(`The "terraform" object does not have an attribute named %q.`, addr.Name),
			Subject:  rng.Ptr(),
		})
		return cty.DynamicVal, diags
	}
}

func nameSuggestion(given string, suggestions []string) string {
	for _, suggestion := range suggestions {
		dist := levenshtein.Distance(given, suggestion, nil)
		if dist < 3 {
			return suggestion
		}
	}
	return ""
}
