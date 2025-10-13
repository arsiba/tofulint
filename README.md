# TofuLint
[![Build Status](https://github.com/arsiba/tofulint/workflows/build/badge.svg?branch=master)](https://github.com/arsiba/tofulint/actions)
[![GitHub release](https://img.shields.io/github/release/arsiba/tofulint.svg)](https://github.com/arsiba/tofulint/releases/latest)
[![Terraform Compatibility](https://img.shields.io/badge/terraform-%3E%3D%201.0-blue)](docs/user-guide/compatibility.md)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsiba/tofulint)](https://goreportcard.com/report/github.com/arsiba/tofulint)
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/tflint.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/tflint)

A **pluggable** [OpenTofu](https://opentofu.org/) linter inspired by TFLint.

## ⚠️ Disclaimer
`TofuLint` is an **experimental** fork of `TFLint` that replaces Terraform internals with **OpenTofu**.  
It is **highly experimental** and **not production-ready**. Use at your own risk.


## Features

TofuLint is a modular framework where each feature is provided via plugins. Key features include:

- Detect potential errors (e.g., invalid instance types) for major cloud providers: AWS, Azure, GCP.  
- Warn about deprecated syntax and unused declarations.  
- Enforce best practices and naming conventions.  

## Installation
Currently, only one installation method is available:

### Bash (Linux)
```bash
curl -s https://raw.githubusercontent.com/arsiba/tofulint/master/install_linux.sh | bash
````

### Verification
At this stage, no releases are verified or signed.

### Docker
A Docker-based installation will be available in a future release.

## Getting Started
TofuLint comes bundled with a [Terraform language ruleset](https://github.com/arsiba/tofulint-ruleset-opentofu), enabling recommended rules by default.

### Enabling the Terraform Plugin
Declare the plugin block in your `.tflint.hcl`:

```hcl
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}
```

More details: [TFLint Terraform Ruleset Configuration](https://github.com/arsiba/tofulint-ruleset-opentofu/blob/main/docs/configuration.md)

### Cloud Provider Plugins

If you use a cloud provider, install the corresponding plugin:

* [AWS](https://github.com/terraform-linters/tflint-ruleset-aws)
* [Azure](https://github.com/terraform-linters/tflint-ruleset-azurerm)
* [GCP](https://github.com/terraform-linters/tflint-ruleset-google)

Other plugins can be added via `.tflint.hcl` and installed with:

```bash
tofulint --init
```

### Example Plugin Configuration

```hcl
plugin "foo" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/org/tflint-ruleset-foo"

  signing_key = <<-KEY
  -----BEGIN PGP PUBLIC KEY BLOCK-----
  ...
  KEY
}
```

For custom rules, create your own plugin or use Rego policies:

* [Writing Plugins](docs/developer-guide/plugins.md)
* [OPA Ruleset](https://github.com/terraform-linters/tflint-ruleset-opa)


## Usage

By default, TofuLint inspects files in the current directory. Example options:

```bash
$ tofulint --help
Usage:
  tofulint --chdir=DIR/--recursive [OPTIONS]

Application Options:
  -v, --version                         Print TofuLint version
      --init                            Install plugins
      --langserver                      Start language server
  -f, --format=[default|json|checkstyle|junit|compact|sarif] Output format
  -c, --config=FILE                     Config file name (default: .tflint.hcl)
      --ignore-module=SOURCE            Ignore module sources
      --enable-rule=RULE_NAME           Enable rules from the command line
      --disable-rule=RULE_NAME          Disable rules from the command line
      --only=RULE_NAME                  Enable only this rule
      --enable-plugin=PLUGIN_NAME       Enable plugins from the command line
      --var-file=FILE                    Terraform variable file
      --var='foo=bar'                    Set a Terraform variable
      --call-module-type=[all|local|none] Types of module to call (default: local)
      --chdir=DIR                        Change working directory
      --recursive                        Run recursively in subdirectories
      --filter=FILE                       Filter issues by file names/globs
      --force                             Return zero exit code even if issues found
      --minimum-failure-severity=[error|warning|notice] Minimum severity for non-zero exit
      --color                             Enable colorized output
      --no-color                          Disable colorized output
      --fix                               Automatically fix issues
      --no-parallel-runners               Disable parallelism

Help Options:
  -h, --help                             Show this help message
```

See [User Guide](docs/user-guide) for more details.

## Debugging

Enable detailed logs using the `TFLINT_LOG` environment variable:

```bash
$ TFLINT_LOG=debug tofulint
```
## Developing

See [Developer Guide](docs/developer-guide) for instructions on contributing and building plugins.

## Security

For reporting security issues, refer to our [security policy](SECURITY.md).

