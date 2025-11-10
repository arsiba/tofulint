package plugin

import (
	"github.com/arsiba/tofulint-plugin-sdk/plugin/host2plugin"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

// PluginRoot is the root directory of the plugins
// This variable is exposed for testing.
var (
	PluginRoot      = "~/.tflint.d/plugins"
	localPluginRoot = "./.tflint.d/plugins"
)

// SDKVersionConstraints is the version constraint of the supported SDK version.
var SDKVersionConstraints = version.MustConstraints(version.NewConstraint(">= 0.0.7"))

// Plugin is an object handling plugins
// Basically, it is a wrapper for go-plugin and provides an API to handle them collectively.
type Plugin struct {
	RuleSets map[string]*host2plugin.Client

	clients map[string]*plugin.Client
}

// Clean is a helper for ending plugin processes
func (p *Plugin) Clean() {
	for _, client := range p.clients {
		client.Kill()
	}
}
