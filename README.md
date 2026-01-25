# ok

A unified task runner that discovers and executes tasks from multiple build tools.

## Overview

`ok` automatically discovers tasks from various build tool configuration files in your project (like `Makefile` and `package.json`) and provides a single interface to run them all. No more remembering whether to use `make build`, `npm run build`, or another tool-specific command.

## Features

- **Automatic task discovery** - Scans your project for supported build tool files
- **Unified interface** - Run tasks from any tool using the same command
- **Multiple tool support** - Currently supports:
  - Make (Makefile)
  - NPM (package.json scripts)
- **Fast concurrent discovery** - Uses goroutines to discover tasks in parallel
- **Simple CLI** - Minimal flags, maximum productivity

## Installation

```bash
go install github.com/broothie/ok@latest
```

## Usage

### List all available tasks

Run `ok` without arguments to see all discovered tasks:

```bash
ok
```

This will show tasks from all supported tools in your project.

### Run a task

```bash
ok <task-name> [args...]
```

For example:
```bash
ok build
ok test --verbose
ok start
```

`ok` will automatically find the task and run it with the appropriate tool, passing along any additional arguments.

### Options

- `-h, --help` - Show command help
- `-d, --directory <path>` - Directory to run command from (default: `.`)
- `-t, --timeout <duration>` - Command timeout (default: `1s`)

### Examples

```bash
# List all tasks in current directory
ok

# Run a build task
ok build

# Run tests with arguments
ok test --coverage

# Run from a different directory
ok -d ./my-project build

# Set a longer timeout
ok -t 5m integration-test
```

## How It Works

1. `ok` scans your project for supported build tool files
2. For each file found, it parses the available tasks
3. When you run a task, `ok` executes it using the appropriate underlying tool
4. All output is streamed directly to your terminal

## Supported Tools

### Make
- Discovers targets from `Makefile`
- Runs tasks using `make <target>`

### NPM
- Discovers scripts from `package.json`
- Runs tasks using `npm run <script>`

## Development

### Requirements

- Go 1.25.6 or later

### Build

```bash
go build
```

### Test

```bash
go test ./...
```

## License

See LICENSE file for details.
