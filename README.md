# wingetcfg

Go library for programmatically generating Windows Package Manager (WinGet) DSC configuration files.

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Copyright 2026 dragonsecurity

## Installation

```sh
go get eye.dragonsecurity.io/wingetcfg
```

## Usage

```go
package main

import (
	"log"

	"eye.dragonsecurity.io/wingetcfg/wingetcfg"
)

func main() {
	cfg := wingetcfg.NewWingetCfg()

	pkg, err := wingetcfg.InstallPackage("", "install firefox", "Mozilla.Firefox", "", "", true)
	if err != nil {
		log.Fatal(err)
	}

	cfg.AddResource(pkg)

	if err := cfg.WriteConfigFile("config.winget"); err != nil {
		log.Fatal(err)
	}
}
```

See [examples/example.go](examples/example.go) for more detailed usage.

## Supported DSC Resource Types

- **WinGet packages** - install/uninstall via Windows Package Manager
- **MSI packages** - install/uninstall MSI installers
- **Registry keys and values** - create and modify Windows registry entries
- **Local user accounts** - manage local Windows users
- **Local groups and membership** - manage local groups and their members
- **PowerShell script execution** - run PowerShell scripts as DSC resources

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
