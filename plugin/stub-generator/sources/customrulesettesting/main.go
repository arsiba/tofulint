package main

import (
	"github.com/arsiba/tofulint/plugin/stub-generator/sources/customrulesettesting/custom"
	"github.com/arsiba/tofulint/plugin/stub-generator/sources/customrulesettesting/rules"
	"github.com/arsiba/tofulint-plugin-sdk/plugin"
	"github.com/arsiba/tofulint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &custom.RuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "customrulesettesting",
				Version: "0.1.0",
				Rules: []tflint.Rule{
					rules.NewAwsInstanceExampleTypeRule(),
				},
			},
		},
	})
}
