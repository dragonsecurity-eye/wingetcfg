# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go library for programmatically generating Windows Package Manager (WinGet) DSC (Desired State Configuration) YAML files. Part of the DragonEye project under DragonSecurity. The library produces `.winget` configuration files conforming to the DSC schema v0.2.

Module: `eye.dragonsecurity.io/wingetcfg`

## Build & Test

```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
go test ./wingetcfg/    # Run tests for the core package only
go vet ./...            # Static analysis
```

No test files exist yet. The only external dependency is `gopkg.in/yaml.v3`.

## Architecture

The library lives in the `wingetcfg/` package. Each file wraps a specific DSC resource type:

- **wingetcfg.go** - Core types (`WinGetCfg`, `WinGetResource`, `WinGetProperties`, `WinGetDirectives`), config creation (`NewWingetCfg`), and YAML file output (`WriteConfigFile`). All resource builders return `*WinGetResource` which gets added via `AddResource` (desired state) or `AddAssertion` (precondition checks).
- **winget_package.go** - `Microsoft.WinGet.DSC/WinGetPackage` - install/uninstall winget packages
- **winget_msi.go** - `xPSDesiredStateConfiguration/xMsiPackage` - install/uninstall MSI packages
- **winget_registry.go** - `xPSDesiredStateConfiguration/xRegistry` - add/remove registry keys and values
- **winget_local_user.go** - `xPSDesiredStateConfiguration/xUser` - manage local user accounts
- **winget_local_group.go** - `xPSDesiredStateConfiguration/xGroup` - manage local groups and membership
- **winget_powershell.go** - `openuem/Powershell` - custom resource for executing PowerShell scripts
- **winget_error_codes.go** - Lookup map of WinGet hex error codes to descriptions

### Pattern for adding new resource types

Each resource file follows a consistent pattern:
1. Define the DSC resource constant (e.g., `WinGetPackageResource`)
2. Create convenience functions (e.g., `InstallPackage`/`UninstallPackage`) that call a lower-level `New*Resource` constructor
3. The constructor builds a `WinGetResource` with `Resource`, optional `ID`, `Directives`, and `Settings` map
4. Use `EnsurePresent`/`EnsureAbsent` constants via `SetEnsureValue()` for the ensure pattern

The `examples/` directory contains a standalone example demonstrating package, and registry resource usage.