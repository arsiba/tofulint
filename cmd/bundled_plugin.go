package cmd

import (
	"fmt"

	"github.com/arsiba/tofulint-plugin-sdk/plugin"
	"github.com/arsiba/tofulint-plugin-sdk/tflint"
	"github.com/arsiba/tofulint-ruleset-opentofu/project"
	"github.com/arsiba/tofulint-ruleset-opentofu/rules"
	"github.com/arsiba/tofulint-ruleset-opentofu/terraform"
)

func (cli *CLI) actAsBundledPlugin() int {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &terraform.RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "opentofu",
				Version: fmt.Sprintf("%s-bundled", project.Version),
			},
			PresetRules: rules.PresetRules,
		},
	})
	return ExitCodeOK
}
