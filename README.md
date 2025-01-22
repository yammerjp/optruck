# optruck

optruck is a command-line tool for managing secrets between 1Password Vault and your applications, with support for template-based restoration.

## Status

This project is currently in the early stages of development.

## Installation

```
go install github.com/yammerjp/optruck
```

## Usage

```
Usage: optruck <command> [flags]

A command-line tool for managing secrets between 1Password Vault and your applications, with support for template-based restoration.

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  upload --item=STRING [flags]
    Upload secrets to a 1Password Vault.

  template --item=STRING [flags]
    Generate a restoration template file.

  mirror --item=STRING [flags]
    Upload secrets and generate a template in one step.

  version [flags]
    Show version information.

Run "optruck <command> --help" for more information on a command.
```
