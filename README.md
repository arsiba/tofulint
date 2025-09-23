# TofuLint
[![Build Status](https://github.com/arsiba/tofulint/workflows/build/badge.svg?branch=master)](https://github.com/arsiba/tofulint/actions)
[![GitHub release](https://img.shields.io/github/release/terraform-linters/tflint.svg)](https://github.com/arsiba/tofulint/releases/latest)
[![Terraform Compatibility](https://img.shields.io/badge/terraform-%3E%3D%201.0-blue)](docs/user-guide/compatibility.md)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/arsiba/tofulint)](https://goreportcard.com/report/github.com/arsiba/tofulint)
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/tflint.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/tflint)

A Pluggable [OpenTofu](https://opentofu.org/) Linter

## Features

TofuLint is a framework, forked and modified from TFLint, and each feature is provided by plugins, the key features are as follows:

- Find possible errors (like invalid instance types) for Major Cloud providers (AWS/Azure/GCP).
- Warn about deprecated syntax, unused declarations.
- Enforce best practices, naming conventions.

## Installation

Only one installation method is in this early development stage available:

Bash script (Linux):

```console
curl -s https://raw.githubusercontent.com/arsiba/tofulint/master/install_linux.sh | bash
```

### Verification

Right now, no releases are verified or signed.

### Docker

In a future release, Instead of installing directly, you will be able to use an Docker image:

## Getting Started

First, enable rules for [Terraform Language](https://www.terraform.io/language) (e.g. warn about deprecated syntax, unused declarations). [TFLint Ruleset for Terraform Language](https://github.com/terraform-linters/tflint-ruleset-terraform) is bundled with TFLint, so you can use it without installing it separately.

The bundled plugin enables the "recommended" preset by default, but you can disable the plugin or use a different preset. Declare the plugin block in `.tflint.hcl` like this:

```hcl
plugin "terraform" {
  enabled = true
  preset  = "recommended"
}
```

See the [tflint-ruleset-terraform documentation](https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/configuration.md) for more information.

Next, If you are using an AWS/Azure/GCP provider, it is a good idea to install the plugin and try it according to each usage:

- [Amazon Web Services](https://github.com/terraform-linters/tflint-ruleset-aws)
- [Microsoft Azure](https://github.com/terraform-linters/tflint-ruleset-azurerm)
- [Google Cloud Platform](https://github.com/terraform-linters/tflint-ruleset-google)

If you want to extend TFLint with other plugins, you can declare the plugins in the config file and easily install them with `tflint --init`.

```hcl
plugin "foo" {
  enabled = true
  version = "0.1.0"
  source  = "github.com/org/tflint-ruleset-foo"

  signing_key = <<-KEY
  -----BEGIN PGP PUBLIC KEY BLOCK-----

  mQINBFzpPOMBEADOat4P4z0jvXaYdhfy+UcGivb2XYgGSPQycTgeW1YuGLYdfrwz
  9okJj9pMMWgt/HpW8WrJOLv7fGecFT3eIVGDOzyT8j2GIRJdXjv8ZbZIn1Q+1V72
  AkqlyThflWOZf8GFrOw+UAR1OASzR00EDxC9BqWtW5YZYfwFUQnmhxU+9Cd92e6i
  ...
  KEY
}
```

See also [Configuring Plugins](docs/user-guide/plugins.md).

If you want to add custom rules that are not in existing plugins, you can build your own plugin or write your own policy in Rego. See [Writing Plugins](docs/developer-guide/plugins.md) or [OPA Ruleset](https://github.com/terraform-linters/tflint-ruleset-opa).

## Usage

TofuLint inspects files under the current directory by default. You can change the behavior with the following options/arguments:

```
$ tofulint --help
Usage:
  tofulint --chdir=DIR/--recursive [OPTIONS]

Application Options:
  -v, --version                                                 Print TofuLint version
      --init                                                    Install plugins
      --langserver                                              Start language server
  -f, --format=[default|json|checkstyle|junit|compact|sarif]    Output format
  -c, --config=FILE                                             Config file name (default: .tflint.hcl)
      --ignore-module=SOURCE                                    Ignore module sources
      --enable-rule=RULE_NAME                                   Enable rules from the command line
      --disable-rule=RULE_NAME                                  Disable rules from the command line
      --only=RULE_NAME                                          Enable only this rule, disabling all other defaults. Can be specified multiple times
      --enable-plugin=PLUGIN_NAME                               Enable plugins from the command line
      --var-file=FILE                                           Terraform variable file name
      --var='foo=bar'                                           Set a Terraform variable
      --call-module-type=[all|local|none]                       Types of module to call (default: local)
      --chdir=DIR                                               Switch to a different working directory before executing the command
      --recursive                                               Run command in each directory recursively
      --filter=FILE                                             Filter issues by file names or globs
      --force                                                   Return zero exit status even if issues found
      --minimum-failure-severity=[error|warning|notice]         Sets minimum severity level for exiting with a non-zero error code
      --color                                                   Enable colorized output
      --no-color                                                Disable colorized output
      --fix                                                     Fix issues automatically
      --no-parallel-runners                                     Disable per-runner parallelism

Help Options:
  -h, --help                                                    Show this help message
```

See [User Guide](docs/user-guide) for details.

## Debugging

If you don't get the expected behavior, you can see the detailed logs when running with `TFLINT_LOG` environment variable.

```console
$ TFLINT_LOG=debug tflint
```

## Developing

See [Developer Guide](docs/developer-guide).

## Security

Please refer our [security policy](SECURITY.md).